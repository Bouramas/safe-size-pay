package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"safe-size-pay/cmd/resources"
	"safe-size-pay/internal/constants"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func AuthHandler() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			jwtToken := r.Header.Get("Authorization")
			if jwtToken == "" {
				writeJSONUnauthorizedError(w, "Unauthorized access - missing Authorization header")
				return
			}

			if len(jwtToken) > 7 && jwtToken[:7] == "Bearer " {
				jwtToken = jwtToken[7:]
			} else {
				writeJSONUnauthorizedError(w, "Invalid authorization header scheme.")
				return
			}

			if jwtToken == "" {
				writeJSONUnauthorizedError(w, "No authorization token provided.")
				return
			}

			claims := &resources.Claims{}
			token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
				// Ensure the signing method is HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				// Return the secret key used for signing
				return []byte(constants.SecretKey), nil
			})

			if err != nil || !token.Valid {
				writeJSONUnauthorizedError(w, "Unauthorized access - Invalid Token")
				return
			}

			if strings.TrimSpace(claims.UserID) == "" {
				writeJSONUnauthorizedError(w, "Unauthorized access - Invalid Token - Missing UserID")
				return
			}

			ctx := context.WithValue(r.Context(), constants.CtxTokenKey, token)
			ctx = context.WithValue(ctx, constants.CtxClaimsKey, claims)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func writeJSONUnauthorizedError(w http.ResponseWriter, msg string) {
	errObject := map[string]interface{}{"error": true, "code": http.StatusUnauthorized, "message": msg}
	res, _ := json.Marshal(errObject)
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write(res)
}
