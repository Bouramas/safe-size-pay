package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"safe-size-pay/cmd/resources"
	"safe-size-pay/internal/constants"

	"golang.org/x/crypto/bcrypt"
)

type DBService struct {
	DB *sql.DB
}

func NewDBService(db *sql.DB) *DBService {
	cs := new(DBService)
	cs.DB = db
	return cs
}

// CreateUser creates a new User.
func (db *DBService) CreateUser(u *resources.User) error {
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	query := "INSERT INTO Users (id, name, email, password_hash) VALUES (UUID_TO_BIN(?, true), ?, ?, ?)"
	_, err = db.DB.Exec(query, u.ID, u.Name, u.Email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("error creating new User: %v", err)
	}
	return nil
}

// ValidateLogin - checks whether the specified email and password are matched with a user
func (db *DBService) ValidateLogin(email, plainPassword string) (*resources.User, error) {
	var u resources.User
	query := "SELECT id, name, email, password_hash FROM Users WHERE email = ?"
	err := db.DB.QueryRow(query, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}

	// Compare hashed password with the provided one
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword)) == nil {
		return nil, fmt.Errorf("invalid password: %v", err)
	}

	return &u, nil
}

// GetUserByID retrieves a User by ID.
func (db *DBService) GetUserByID(userID int) (*resources.User, error) {
	var u resources.User
	query := "SELECT BIN_TO_UUID(id, true), name, email FROM Users WHERE id = UUID_TO_BIN(?, true)"
	err := db.DB.QueryRow(query, userID).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return nil, fmt.Errorf("error retrieving User by ID: %v", err)
	}
	return &u, nil
}

// GetUserByEmail retrieves a User by his/her email.
func (db *DBService) GetUserByEmail(email string) (*resources.User, error) {
	var u resources.User
	query := "SELECT BIN_TO_UUID(id, true), name, email, password_hash FROM Users WHERE email = ?"
	err := db.DB.QueryRow(query, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return nil, fmt.Errorf("error retrieving User by Email: %v", err)
	}
	return &u, nil
}

// CreateTransaction adds a new Transaction to the database.
func (db *DBService) CreateTransaction(description, userID string, amount float64) (*resources.Transaction, error) {
	createdAt := time.Now()
	id := uuid.New().String()
	query := "INSERT INTO Transactions (id, user_id, description, amount, created_at) VALUES (UUID_TO_BIN(?, true),UUID_TO_BIN(?, true),?, ?, ?)"
	_, err := db.DB.Exec(query, id, userID, description, amount, createdAt)
	if err != nil {
		return nil, fmt.Errorf("error creating new Transaction: %v", err)
	}

	return &resources.Transaction{
		ID:          id,
		UserID:      userID,
		OrderID:     nil,
		Description: description,
		Amount:      amount,
		OrderStatus: constants.OrderStatusPending,
		CreatedAt:   createdAt,
	}, nil
}

// GetTransactions retrieves transactions of the specified user.
func (db *DBService) GetTransactions(userID string) ([]*resources.Transaction, error) {

	query := "SELECT BIN_TO_UUID(id, true), BIN_TO_UUID(user_id, true), order_id, COALESCE(order_msg, ''), description, amount, order_status, created_at, updated_at FROM Transactions WHERE user_id=UUID_TO_BIN(?, true)"
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transactions: %v", err)
	}

	defer rows.Close()
	transactions := make([]*resources.Transaction, 0)
	for rows.Next() {
		var t resources.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.OrderID, &t.OrderMsg, &t.Description, &t.Amount, &t.OrderStatus, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning transaction: %v", err)
		}
		transactions = append(transactions, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with result set: %v", err)
	}

	return transactions, nil
}

// MarkTransactionFailed marks an existing transaction as failed
func (db *DBService) MarkTransactionFailed(id, msg string) error {
	query := "UPDATE Transactions SET order_status=?, order_msg=? WHERE id=UUID_TO_BIN(?, true)"
	result, err := db.DB.Exec(query, constants.OrderStatusFailed, msg, id)
	if err != nil {
		return fmt.Errorf("MarkTransactionFailed | error updating Transaction: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// PatchTransactionOrderID updates an existing transaction, setting the order ID and status
func (db *DBService) PatchTransactionOrderID(id string, orderID int) error {
	query := "UPDATE Transactions SET order_id=? WHERE id=UUID_TO_BIN(?, true)"
	result, err := db.DB.Exec(query, orderID, id)
	if err != nil {
		return fmt.Errorf("PatchTransactionOrderID | error updating Transaction: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// PatchTransactionSuccess updates an existing transaction, setting the order ID and status
func (db *DBService) PatchTransactionSuccess(id string) error {
	query := "UPDATE Transactions SET order_status=? WHERE id=UUID_TO_BIN(?, true)"
	result, err := db.DB.Exec(query, constants.OrderStatusSuccess, id)
	if err != nil {
		return fmt.Errorf("PatchTransactionSuccess | error updating Transaction: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteTransactionByID deletes a transaction by its ID.
func (db *DBService) DeleteTransactionByID(id string) error {
	query := "DELETE FROM Transactions WHERE id=UUID_TO_BIN(?, true)"
	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting transaction: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
