package server

import "crypto/rsa"

// Типы сохраняемых данных
const (
	TextDataType     = "text"
	FileDataType     = "file"
	CardDataType     = "card"
	PasswordDataType = "password"
)

// ServerOptions содержит настройки серверной части приложения.
var ServerOptions struct {
	DBaddress        string
	PrivateKeyPath   string
	PrivateCryptoKey *rsa.PrivateKey
	GrpcHost         string
}

// SetConfig устанавливает значение настроек по умолчанию.
func SetConfig() {
	ServerOptions.DBaddress = "postgres://postgres:postgres1@localhost:5432/metrics?sslmode=disable"
	ServerOptions.GrpcHost = "localhost:8080"
	ServerOptions.PrivateKeyPath = "configs/server/private.key"
}
