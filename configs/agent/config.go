package agent

import "crypto/rsa"

// Типы сохраняемых данных
const (
	TextDataType     = "text"
	FileDataType     = "file"
	CardDataType     = "card"
	PasswordDataType = "password"
)

// Команды
const (
	CRegestration  = "reg"
	CAuthorization = "log"
	CPasswordSave  = "pass_save"
	CPasswordGet   = "pass_get"
	CTextSave      = "text_save"
	CTextGet       = "text_get"
	CCardSave      = "card_save"
	CCardGet       = "card_get"
	CFileSave      = "file_save"
	CFileGet       = "file_get"
)

// Режимы обработки данных
const (
	GetMode = "get"
	AddMode = "add"
)

// ClientOptions содержит настройки клиентской части приложения.
var ClientOptions struct {
	PublicKeyPath   string
	PrivateKeyPath  string
	PublicCryptoKey *rsa.PublicKey
	GrpcHost        string
	Command         string
	Data            string
	FileStoragePath string
}

// SetConfig устанавливает значение настроек по умолчанию.
func SetConfig() {
	ClientOptions.PublicKeyPath = "configs/agent/public.key"
	ClientOptions.PrivateKeyPath = "configs/server/private.key"
	ClientOptions.GrpcHost = "localhost:8080"
	ClientOptions.FileStoragePath = "outcome_files"
}
