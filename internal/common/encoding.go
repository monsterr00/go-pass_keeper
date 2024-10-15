package applayer

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/monsterr00/go-pass-keeper/internal/models"
)

// EncodeToBytes приводит данных к типу []byte
func EncodeToBytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

// DecodeToText переводит данные из []byte в models.Text
func DecodeToText(s []byte) models.Text {
	t := models.Text{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&t)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

// DecodeToCard переводит данные из []byte в models.Card
func DecodeToCard(s []byte) models.Card {
	c := models.Card{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// DecodeToPassword переводит данные из []byte в models.Password
func DecodeToPassword(s []byte) models.Password {
	p := models.Password{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

// DecodeToFile переводит данные из []byte в models.File
func DecodeToFile(s []byte) models.File {
	f := models.File{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&f)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
