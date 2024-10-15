package models

import "time"

type Card struct {
	UserLogin  string
	Title      string
	Bank       string
	CardNumber string
	CVV        string
	DateExpire time.Time
	CardHolder string
	CreatedAt  time.Time
}
