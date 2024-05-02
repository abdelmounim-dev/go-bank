package main

import "math/rand"

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func NewAccount(info CreateAccountRequest) *Account {
	return &Account{
		FirstName: info.FirstName,
		LastName:  info.LastName,
		Number:    int64(rand.Intn(1e6)),
	}
}
