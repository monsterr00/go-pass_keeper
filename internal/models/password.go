package models

import "time"

type Password struct {
	UserLogin string
	Title     string
	Login     string
	Password  string
	CreatedAt time.Time
}
