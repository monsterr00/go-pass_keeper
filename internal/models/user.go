package models

import "time"

type User struct {
	Login     string
	Password  string
	CreatedAt time.Time
}
