package applayer

import (
	"context"

	config "github.com/monsterr00/go-pass-keeper/configs/server"
	"github.com/monsterr00/go-pass-keeper/internal/models"
	"github.com/monsterr00/go-pass-keeper/internal/server/storelayer"
)

type app struct {
	store storelayer.Store
}

type App interface {
	CloseDB() error
	CreateText(ctx context.Context, text models.Text) error
	UpdateText(ctx context.Context, text models.Text) error
	CheckTextDuplicate(ctx context.Context, login string, title string) error
	GetText(ctx context.Context, login string, title string) (models.Text, error)
	CreatePassword(ctx context.Context, password models.Password) error
	UpdatePassword(ctx context.Context, password models.Password) error
	CheckPasswordDuplicate(ctx context.Context, login string, title string) error
	GetPassword(ctx context.Context, login string, title string) (models.Password, error)
	CreateCard(ctx context.Context, card models.Card) error
	UpdateCard(ctx context.Context, card models.Card) error
	CheckCardDuplicate(ctx context.Context, login string, title string) error
	GetCard(ctx context.Context, login string, title string) (models.Card, error)
	CreateFile(ctx context.Context, file models.File) error
	UpdateFile(ctx context.Context, file models.File) error
	CheckFileDuplicate(ctx context.Context, login string, title string) error
	GetFile(ctx context.Context, login string, title string) (models.File, error)
	Add(ctx context.Context, dataType string, data []byte) error
	Get(ctx context.Context, dataType string, data []byte) ([]byte, error)
}

// New инициализирует уровень app.
func New(storeLayer storelayer.Store) *app {
	return &app{
		store: storeLayer,
	}
}

// CloseDB закрывает соединения к БД.
func (api *app) CloseDB() error {
	err := api.store.Close()
	if err != nil {
		return err
	}
	return nil
}

// CreateText создает новую запись с текстом в БД.
func (api *app) CreateText(ctx context.Context, text models.Text) error {
	return api.store.TextCreate(ctx, text)
}

// UpdateText обновляет данные о текстовой записи в БД.
func (api *app) UpdateText(ctx context.Context, text models.Text) error {
	return api.store.TextUpdate(ctx, text)
}

// CheckTextDuplicate проверяет наличие текстовой записи в БД.
func (api *app) CheckTextDuplicate(ctx context.Context, login string, title string) error {
	return api.store.CheckDuplicateText(ctx, login, title)
}

// GetText возвращает текстовую запись из БД.
func (api *app) GetText(ctx context.Context, login string, title string) (models.Text, error) {
	return api.store.TextGet(ctx, login, title)
}

// CreatePassword создает новую запись с логином/паролем в БД.
func (api *app) CreatePassword(ctx context.Context, password models.Password) error {
	return api.store.PasswordCreate(ctx, password)
}

// UpdatePassword обновляет данные о записи логина/пароля в БД.
func (api *app) UpdatePassword(ctx context.Context, password models.Password) error {
	return api.store.PasswordUpdate(ctx, password)
}

// CheckPasswordDuplicate проверяет наличие записи логина/пароля в БД.
func (api *app) CheckPasswordDuplicate(ctx context.Context, login string, title string) error {
	return api.store.CheckDuplicatePassword(ctx, login, title)
}

// GetPassword возвращает запись с логином/паролем из БД.
func (api *app) GetPassword(ctx context.Context, login string, title string) (models.Password, error) {
	return api.store.PasswordGet(ctx, login, title)
}

// CreateCard создает новую запись с данными банковской карты в БД.
func (api *app) CreateCard(ctx context.Context, card models.Card) error {
	return api.store.CardCreate(ctx, card)
}

// UpdateCard обновляет данные о записи с данными банковской карты в БД.
func (api *app) UpdateCard(ctx context.Context, card models.Card) error {
	return api.store.CardUpdate(ctx, card)
}

// CheckCardDuplicate проверяет наличие записи с данными банковской карты в БД.
func (api *app) CheckCardDuplicate(ctx context.Context, login string, title string) error {
	return api.store.CheckDuplicateCard(ctx, login, title)
}

// GetCard возвращает запись с данными банковской карты из БД.
func (api *app) GetCard(ctx context.Context, login string, title string) (models.Card, error) {
	return api.store.CardGet(ctx, login, title)
}

// CreateFile создает новую запись с содержимым файла в БД.
func (api *app) CreateFile(ctx context.Context, file models.File) error {
	return api.store.FileCreate(ctx, file)
}

// UpdateFile обновляет данные с содержимым файла в БД.
func (api *app) UpdateFile(ctx context.Context, file models.File) error {
	return api.store.FileUpdate(ctx, file)
}

// CheckFileDuplicate проверяет наличие записи с содержимым файла в БД.
func (api *app) CheckFileDuplicate(ctx context.Context, login string, title string) error {
	return api.store.CheckDuplicateFile(ctx, login, title)
}

// GetCard возвращает запись с данными банковской карты из БД.
func (api *app) GetFile(ctx context.Context, login string, title string) (models.File, error) {
	return api.store.FileGet(ctx, login, title)
}

// Add инициализирует стратегию в зависимости от типа сохраняемых данных и передает данные для сохранения.
func (api *app) Add(ctx context.Context, dataType string, data []byte) error {
	dataStrat := initData(api)

	switch dataType {
	case config.TextDataType:
		text := &text{}
		dataStrat.setDataAlgo(text)
	case config.CardDataType:
		card := &card{}
		dataStrat.setDataAlgo(card)
	case config.FileDataType:
		file := &file{}
		dataStrat.setDataAlgo(file)
	case config.PasswordDataType:
		password := &password{}
		dataStrat.setDataAlgo(password)
	}

	dataStrat.create(ctx, data)

	return nil
}

// Get инициализирует стратегию в зависимости от типа сохраняемых данных и возвращает данные из БД.
func (api *app) Get(ctx context.Context, dataType string, data []byte) ([]byte, error) {
	dataStrat := initData(api)

	switch dataType {
	case config.TextDataType:
		text := &text{}
		dataStrat.setDataAlgo(text)
	case config.CardDataType:
		card := &card{}
		dataStrat.setDataAlgo(card)
	case config.FileDataType:
		file := &file{}
		dataStrat.setDataAlgo(file)
	case config.PasswordDataType:
		password := &password{}
		dataStrat.setDataAlgo(password)
	}

	return dataStrat.get(ctx, data)
}
