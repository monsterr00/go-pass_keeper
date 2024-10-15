package applayer

import (
	"context"
	"log"

	e "github.com/monsterr00/go-pass-keeper/internal/common"
)

type dataAlgo interface {
	create(ctx context.Context, data *dataStrat)
	get(ctx context.Context, data *dataStrat) ([]byte, error)
}

type dataStrat struct {
	dataBytes []byte
	dataAlgo  dataAlgo
	app       App
}

func initData(appLayer App) *dataStrat {
	return &dataStrat{
		app: appLayer,
	}
}

func (d *dataStrat) setDataAlgo(a dataAlgo) {
	d.dataAlgo = a
}

func (d *dataStrat) create(ctx context.Context, data []byte) {
	d.dataBytes = data
	d.dataAlgo.create(ctx, d)
}

func (d *dataStrat) get(ctx context.Context, data []byte) ([]byte, error) {
	d.dataBytes = data
	return d.dataAlgo.get(ctx, d)
}

type file struct{}
type card struct{}
type text struct{}
type password struct{}

func (t *file) create(ctx context.Context, d *dataStrat) {
	var err error
	var isDuplicate = false

	file := e.DecodeToFile(d.dataBytes)

	err = d.app.CheckFileDuplicate(ctx, file.UserLogin, file.Title)
	if err != nil {
		isDuplicate = true
	}

	switch isDuplicate {
	case true:
		err = d.app.UpdateFile(ctx, file)
	case false:
		err = d.app.CreateFile(ctx, file)
	}

	if err != nil {
		log.Println(err)
	}
}

func (t *file) get(ctx context.Context, d *dataStrat) ([]byte, error) {
	incFile := e.DecodeToFile(d.dataBytes)

	file, err := d.app.GetFile(ctx, incFile.UserLogin, incFile.Title)
	if err != nil {
		log.Println(err)
	}

	return e.EncodeToBytes(file), nil
}

func (t *card) create(ctx context.Context, d *dataStrat) {
	var err error
	var isDuplicate = false

	card := e.DecodeToCard(d.dataBytes)

	err = d.app.CheckCardDuplicate(ctx, card.UserLogin, card.Title)
	if err != nil {
		isDuplicate = true
	}

	switch isDuplicate {
	case true:
		err = d.app.UpdateCard(ctx, card)
	case false:
		err = d.app.CreateCard(ctx, card)
	}

	if err != nil {
		log.Println(err)
	}
}

func (t *card) get(ctx context.Context, d *dataStrat) ([]byte, error) {
	incCard := e.DecodeToCard(d.dataBytes)

	card, err := d.app.GetCard(ctx, incCard.UserLogin, incCard.Title)
	if err != nil {
		log.Println(err)
	}

	return e.EncodeToBytes(card), nil
}

func (t *text) create(ctx context.Context, d *dataStrat) {
	var err error
	var isDuplicate = false

	text := e.DecodeToText(d.dataBytes)

	err = d.app.CheckTextDuplicate(ctx, text.UserLogin, text.Title)
	if err != nil {
		isDuplicate = true
	}

	switch isDuplicate {
	case true:
		err = d.app.UpdateText(ctx, text)
	case false:
		err = d.app.CreateText(ctx, text)
	}

	if err != nil {
		log.Println(err)
	}
}

func (t *text) get(ctx context.Context, d *dataStrat) ([]byte, error) {
	incText := e.DecodeToText(d.dataBytes)

	text, err := d.app.GetText(ctx, incText.UserLogin, incText.Title)
	if err != nil {
		log.Println(err)
	}

	return e.EncodeToBytes(text), nil
}

func (t *password) create(ctx context.Context, d *dataStrat) {
	var err error
	var isDuplicate = false

	password := e.DecodeToPassword(d.dataBytes)

	err = d.app.CheckPasswordDuplicate(ctx, password.UserLogin, password.Title)
	if err != nil {
		isDuplicate = true
	}

	switch isDuplicate {
	case true:
		err = d.app.UpdatePassword(ctx, password)
	case false:
		err = d.app.CreatePassword(ctx, password)
	}

	if err != nil {
		log.Println(err)
	}
}

func (t *password) get(ctx context.Context, d *dataStrat) ([]byte, error) {
	incPassword := e.DecodeToPassword(d.dataBytes)

	password, err := d.app.GetPassword(ctx, incPassword.UserLogin, incPassword.Title)
	if err != nil {
		log.Println(err)
	}

	return e.EncodeToBytes(password), nil
}
