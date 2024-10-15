package models

import "time"

type File struct {
	UserLogin string
	Title     string
	FileName  string
	File      []byte
	DataType  string
	CreatedAt time.Time
}
