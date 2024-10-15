package applayer

import (
	"io"
	"log"
	"os"
	"strings"

	config "github.com/monsterr00/go-pass-keeper/configs/agent"
	"github.com/monsterr00/go-pass-keeper/internal/agent/storelayer"
)

type app struct {
	store storelayer.Store
}

type App interface {
	ParseData() []string
	ReadFromFile(path string) []byte
	WriteToFile(s []byte, file string)
}

func New(storeLayer storelayer.Store) *app {
	return &app{
		store: storeLayer,
	}
}

// ParseData разбивает входящую строку на части по разделяющему символу ";"
func (api *app) ParseData() []string {
	return strings.FieldsFunc(config.ClientOptions.Data, func(r rune) bool {
		return r == ';'
	})
}

// ReadFromFile считывает содержимое файла и переводит в []byte
func (api *app) ReadFromFile(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

// WriteToFile сохраняет данные типа []byte в файл
func (api *app) WriteToFile(s []byte, file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	f.Write(s)
}
