package httplayer

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/monsterr00/go-pass-keeper/configs/agent"
	"github.com/monsterr00/go-pass-keeper/internal/agent/applayer"
	e "github.com/monsterr00/go-pass-keeper/internal/common"
	"github.com/monsterr00/go-pass-keeper/internal/models"
	pb "github.com/monsterr00/go-pass-keeper/internal/server/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type httpAPI struct {
	app        applayer.App
	grpcConn   *grpc.ClientConn
	grpcClient pb.KeeperDataClient
}

// New инициализирует уровень app
func New(appLayer applayer.App) *httpAPI {
	api := &httpAPI{
		app:        appLayer,
		grpcConn:   nil,
		grpcClient: nil,
	}

	return api
}

// Engage запускает grpc-сервер другие службы приложения.
func (api *httpAPI) Engage() {
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		api.stopServer()
		close(idleConnsClosed)
	}()

	api.startServer()
	api.startServices()

	<-idleConnsClosed
	log.Printf("Client Shutdown gracefully start")
	api.stopServices()
	log.Printf("Client Shutdown gracefully end")
}

// generateCryptoKeys загружает ключи шифрования из файла или генерирует их
func (api *httpAPI) generateCryptoKeys() {
	filePub, err := os.OpenFile(config.ClientOptions.PublicKeyPath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer filePub.Close()

	scanner := bufio.NewScanner(filePub)
	scanner.Scan()
	publicKeyFile := scanner.Bytes()

	if len(publicKeyFile) > 0 {
		publicKeyFile, err = os.ReadFile(config.ClientOptions.PublicKeyPath)
		if err != nil {
			log.Fatal(err)
		}

		publicKeyBlock, _ := pem.Decode(publicKeyFile)
		publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
		if err != nil {
			log.Fatal(err)
		}

		config.ClientOptions.PublicCryptoKey = publicKey

	} else {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatal(err)
		}

		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		})

		filePriv, err := os.OpenFile(config.ClientOptions.PrivateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}

		writer := bufio.NewWriter(filePriv)

		if _, err := writer.Write(privateKeyPEM); err != nil {
			log.Fatal(err)
		}
		if err := writer.Flush(); err != nil {
			log.Fatal(err)
		}

		err = filePriv.Close()
		if err != nil {
			log.Fatal(err)
		}

		publicKey := &privateKey.PublicKey
		config.ClientOptions.PublicCryptoKey = publicKey

		publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
		publicKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		})

		writer = bufio.NewWriter(filePub)
		if _, err := writer.Write(publicKeyPEM); err != nil {
			log.Fatal(err)
		}
		if err := writer.Flush(); err != nil {
			log.Fatal(err)
		}
	}
}

// encrypt используется для шифрования исходящих запросов
func (api *httpAPI) encrypt(body string) string {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, config.ClientOptions.PublicCryptoKey, []byte(body))
	if err != nil {
		log.Fatal(err)
	}
	return string(ciphertext)
}

// stopServices останавливает работу всех сервисов агента
func (api *httpAPI) stopServices() {}

// stopServer останавливает работу серверов
func (api *httpAPI) stopServer() {
	api.grpcConn.Close()
}

// startServer запускает работу серверов
func (api *httpAPI) startServer() {
	api.generateCryptoKeys()

	var err error
	api.grpcConn, err = grpc.NewClient(config.ClientOptions.GrpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	api.grpcClient = pb.NewKeeperDataClient(api.grpcConn)
}

// startServices запускает работу всех сервисов агента
func (api *httpAPI) startServices() {
	api.SendData()
}

// SendData подготавливает и отправляети данные на сервер.
func (api *httpAPI) SendData() {
	var dataGrpc pb.Data
	var mode string

	parsedData := api.app.ParseData()

	switch config.ClientOptions.Command {
	case config.CRegestration:
	case config.CAuthorization:
	case config.CPasswordSave:
		var p models.Password
		p.UserLogin = parsedData[0]
		p.Title = parsedData[1]
		p.Login = parsedData[2]
		p.Password = parsedData[3]

		dataGrpc.Data = e.EncodeToBytes(p)
		dataGrpc.DataType = config.PasswordDataType
		mode = config.AddMode
	case config.CPasswordGet:
		var p models.Password
		p.UserLogin = parsedData[0]
		p.Title = parsedData[1]

		dataGrpc.Data = e.EncodeToBytes(p)
		dataGrpc.DataType = config.PasswordDataType
		mode = config.GetMode
	case config.CTextSave:
		var t models.Text
		t.UserLogin = parsedData[0]
		t.Title = parsedData[1]
		t.Text = parsedData[2]

		dataGrpc.Data = e.EncodeToBytes(t)
		dataGrpc.DataType = config.TextDataType
		mode = config.AddMode
	case config.CTextGet:
		var t models.Text
		t.UserLogin = parsedData[0]
		t.Title = parsedData[1]

		dataGrpc.Data = e.EncodeToBytes(t)
		dataGrpc.DataType = config.TextDataType
		mode = config.GetMode
	case config.CCardSave:
		var c models.Card
		c.UserLogin = parsedData[0]
		c.Title = parsedData[1]
		c.Bank = parsedData[2]
		c.CardNumber = parsedData[3]
		c.CVV = parsedData[4]
		c.CardHolder = parsedData[5]
		c.DateExpire, _ = time.Parse(time.DateOnly, parsedData[6])

		dataGrpc.Data = e.EncodeToBytes(c)
		dataGrpc.DataType = config.CardDataType
		mode = config.AddMode
	case config.CCardGet:
		var c models.Card
		c.UserLogin = parsedData[0]
		c.Title = parsedData[1]

		dataGrpc.Data = e.EncodeToBytes(c)
		dataGrpc.DataType = config.CardDataType
		mode = config.GetMode
	case config.CFileSave:
		var f models.File
		f.UserLogin = parsedData[0]
		f.Title = parsedData[1]
		f.FileName = parsedData[2]
		f.DataType = parsedData[3]
		f.File = api.app.ReadFromFile(parsedData[4])

		dataGrpc.Data = e.EncodeToBytes(f)
		dataGrpc.DataType = config.FileDataType
		mode = config.AddMode
	case config.CFileGet:
		var f models.File
		f.UserLogin = parsedData[0]
		f.Title = parsedData[1]

		dataGrpc.Data = e.EncodeToBytes(f)
		dataGrpc.DataType = config.FileDataType
		mode = config.GetMode
	default:
		log.Println("Unknown command")
		return
	}

	switch mode {
	case config.GetMode:
		resp, err := api.grpcClient.GetData(context.Background(), &pb.GetDataRequest{
			Data: &dataGrpc,
		})

		if err != nil {
			log.Printf("Client: %s\n", err)
		}

		switch dataGrpc.DataType {
		case config.FileDataType:
			file := e.DecodeToFile(resp.Data.Data)
			api.app.WriteToFile(file.File, config.ClientOptions.FileStoragePath+"/"+file.FileName)
		case config.TextDataType:
			fmt.Println(e.DecodeToText(resp.Data.Data))
		case config.CardDataType:
			fmt.Println(e.DecodeToCard(resp.Data.Data))
		case config.PasswordDataType:
			fmt.Println(e.DecodeToPassword(resp.Data.Data))
		}

	case config.AddMode:
		resp, err := api.grpcClient.AddData(context.Background(), &pb.AddDataRequest{
			Data: &dataGrpc,
		})

		if err != nil {
			log.Printf("Client: %s\n", err)
		}
		if resp.Error != "" {
			log.Println(resp.Error)
		}
	}
}
