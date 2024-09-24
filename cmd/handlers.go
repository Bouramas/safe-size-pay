package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"safe-size-pay/cmd/resources"
	"safe-size-pay/cmd/resources/requests"
	"safe-size-pay/cmd/resources/responses"
	"safe-size-pay/internal/constants"
	"safe-size-pay/internal/services/viva"
)

func (s *Server) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !s.hasJsonContentType(r) {
			s.writeJSONError(w, http.StatusBadRequest, "Content-Type is not application/json.")
			return
		}

		loginRequest := &requests.LoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			s.writeJSONError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		err = loginRequest.Validate()
		if err != nil {
			s.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := s.DBService.GetUserByEmail(loginRequest.Email)
		if err != nil || !checkPasswordHash(loginRequest.Password, user.Password) {
			s.writeJSONError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		// Generate JWT Token
		expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
		claims := &resources.Claims{
			UserID: user.ID,
			Email:  user.Email,
			Name:   user.Name,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		var jwtKey = []byte(constants.SecretKey)

		// Create the JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		loginResponse := responses.LoginResponse{
			Token: tokenString,
			Name:  user.Name,
			ID:    user.ID,
		}

		res, _ := json.Marshal(loginResponse)
		_, _ = w.Write(res)
	}
}

func (s *Server) HandlePostSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.Log.Infof("Post Signup")
		if !s.hasJsonContentType(r) {
			s.writeJSONError(w, http.StatusBadRequest, "Content-Type is not application/json.")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		user := &resources.User{}
		if err = json.Unmarshal(body, &user); err != nil {
			s.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		user.ID = uuid.New().String()

		err = user.Validate()
		if err != nil {
			s.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = s.DBService.CreateUser(user)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "duplicate entry") {
				s.writeJSONError(w, http.StatusConflict, "a user with the same email already exists")
				return
			}
			s.writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, _ := json.Marshal(user)
		_, _ = w.Write(res)
	}
}

func (s *Server) HandlePostTransactionHook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.Log.Infof("Post Transactions Hook")
		if !s.hasJsonContentType(r) {
			s.writeJSONError(w, http.StatusBadRequest, "Content-Type is not application/json.")
			return
		}

		// TODO: Parse response and update Transaction in DB accordingly
	}
}

func (s *Server) HandleGetTransactionHook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.Log.Infof("Get Transactions Hook")

		if !s.hasJsonContentType(r) {
			s.writeJSONError(w, http.StatusBadRequest, "Content-Type is not application/json.")
			return
		}

		vivaKey, err := s.VivaClient.GetToken()
		if err != nil {
			s.writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with the key and a custom header
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("test-header", "value")

		s.Log.Infof("Webhook has been validated")


		// TODO: Need to confirm that this request is valid since sending the access token through an api call does not seem valid.
		res, _ := json.Marshal(responses.WebhookResponse{
			Key: vivaKey.AccessToken,
		})
		_, _ = w.Write(res)
	}
}

func (s *Server) HandlePostTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userClaims := ctx.Value(constants.CtxClaimsKey).(*resources.Claims)
		s.Log.Infof("Post Transactions | User ID: %s", userClaims.UserID)

		if !s.hasJsonContentType(r) {
			s.writeJSONError(w, http.StatusBadRequest, "Content-Type is not application/json.")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		trRequest := &requests.TransactionRequest{}
		if err = json.Unmarshal(body, &trRequest); err != nil {
			s.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = trRequest.Validate()
		if err != nil {
			s.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		transaction, err := s.DBService.CreateTransaction(trRequest.Description, userClaims.UserID, trRequest.Amount)
		if err != nil {
			s.writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		orderRequest := &viva.CreateOrderRequest{
			Amount:       viva.AmountConversion(transaction.Amount),
			CustomerTrns: transaction.Description,
			Customer: viva.Customer{
				Email:       userClaims.Email,
				FullName:    userClaims.Name,
				RequestLang: constants.LangEnUS,
			},
		}

		order, err := s.VivaClient.CreateOrder(orderRequest)
		if err != nil {
			s.Log.Errorf("Post Transactions | Error: %s", err.Error())
			_ = s.DBService.MarkTransactionFailed(transaction.ID, fmt.Sprintf("viva failure: %s", err.Error()))
			s.writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = s.DBService.PatchTransactionOrderID(transaction.ID, order.OrderCode)
		if err != nil {
			s.writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, _ := json.Marshal(responses.TransactionResponse{
			RedirectUrl: fmt.Sprintf("%s?ref=%d&color=00BFFF", s.VivaRedirectUrl, order.OrderCode),
			Status:      transaction.OrderStatus,
			ID:          transaction.ID,
		})
		_, _ = w.Write(res)
	}
}

func (s *Server) HandleGetTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userClaims := ctx.Value(constants.CtxClaimsKey).(*resources.Claims)
		s.Log.Infof("Get Transactions | User ID: %s", userClaims.UserID)

		transactions, err := s.DBService.GetTransactions(userClaims.UserID)
		if err != nil {
			s.writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, _ := json.Marshal(transactions)
		_, _ = w.Write(res)
	}
}

func (s *Server) HandleDeleteTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userClaims := ctx.Value(constants.CtxClaimsKey).(*resources.Claims)
		s.Log.Infof("Delete Transactions | User ID: %s", userClaims.UserID)

		id := mux.Vars(r)["id"]

		if err := s.DBService.DeleteTransactionByID(id); err != nil {
			if err == sql.ErrNoRows {
				s.writeJSONError(w, http.StatusNotFound, "Transaction not found")
				return
			}
			s.writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// checkPasswordHash compares the hashed password with the provided password.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
