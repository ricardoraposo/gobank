package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAllAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    number SERIAL,
    balance SERIAL,
    created_at timestamp
    )`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := ` INSERT INTO account
    (first_name, last_name, number, balance, created_at)
        VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
    _, err := s.db.Query("DELETE FROM account WHERE id = $1", id)
	return err
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
    rows, err := s.db.Query("SELECT * FROM account WHERE id = $1", id)
    if err != nil {
        return nil, err
    }
    for rows.Next() {
        return scanIntoAccount(rows)
    }
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAllAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
	return account, err
}
