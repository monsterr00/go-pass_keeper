package models

import "time"

type Text struct {
	UserLogin string
	Title     string
	Text      string
	CreatedAt time.Time
}
