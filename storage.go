package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() (map[int]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgressStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=pw sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	s := &PostgresStore{
		db: db,
	}

	return s, s.Init()
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := "INSERT INTO accounts (first_name, last_name, number, balance) VALUES ($1, $2, $3, $4) returning id"

	err := s.db.QueryRow(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance).Scan(&acc.ID)
	if err != nil {
		return err
	}
	return err
}
func (s *PostgresStore) UpdateAccount(acc *Account) error {
	query := "UPDATE accounts SET first_name=$1, last_name=$2, number=$3, balance=$4 WHERE id=$5"
	_, err := s.db.Exec(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.ID)
	return err
}
func (s *PostgresStore) DeleteAccount(id int) error {
	query := "DELETE FROM accounts WHERE id=$1"
	_, err := s.db.Exec(query, id)
	return err
}
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := "SELECT id, first_name, last_name, number, balance FROM accounts WHERE id=$1"
	result, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	acc := new(Account)
	if result.Next() {
		if err = result.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance); err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("element not found")
}
func (s *PostgresStore) GetAccounts() (map[int]*Account, error) {
	query := "SELECT id, first_name, last_name, number, balance FROM accounts"
	result, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	accounts := map[int]*Account{}
	for result.Next() {
		acc := new(Account)
		if err = result.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance); err != nil {
			return nil, err
		}
		accounts[acc.ID] = acc
	}
	return accounts, err
}

func (s *PostgresStore) Init() error {
	queries := []string{
		"CREATE TABLE IF NOT EXISTS accounts (id SERIAL PRIMARY KEY, first_name varchar(50), last_name varchar(50), number varchar(50), balance varchar(50))",
	}

	for _, query := range queries {
		_, err := s.db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
