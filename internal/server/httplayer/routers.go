package httplayer

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	config "github.com/monsterr00/go-pass-keeper/configs/server"
	"github.com/monsterr00/go-pass-keeper/internal/server/applayer"
	pb "github.com/monsterr00/go-pass-keeper/internal/server/grpc"
	"google.golang.org/grpc"
)

type httpAPI struct {
	app  applayer.App
	grpc *grpc.Server
}

// New инициализирует http-сервер и другие службы приложения.
func New(appLayer applayer.App) *httpAPI {
	return &httpAPI{
		app: appLayer,
		//grpc: grpc.NewServer(grpc.UnaryInterceptor(checkAuthInterceptor)),
		grpc: grpc.NewServer(),
	}
}

// Engage запускает grpc-сервер и другие службы приложения.
func (api *httpAPI) Engage() {
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		api.stopServer()
		close(idleConnsClosed)
	}()

	api.startServices()
	api.startServer()

	<-idleConnsClosed
	log.Printf("Server Shutdown gracefully start")
	api.stopServices()
	log.Printf("Server Shutdown gracefully end")
}

// startServer запускает grpc-сервер
func (api *httpAPI) startServer() {
	api.startGrpc()

	listen, err := net.Listen("tcp", config.ServerOptions.GrpcHost)
	if err != nil {
		log.Fatal(err)
	}

	if err := api.grpc.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

// stopServer останавливает grpc-сервер
func (api *httpAPI) stopServer() {
	api.grpc.Stop()
}

// stopServer останавливает сервисы
func (api *httpAPI) stopServices() {
	api.closeDB()
}

// startServices запускает сервисы
func (api *httpAPI) startServices() {
	api.generateCryptoKeys()
}

// generateCryptoKeys загружает приватный ключ из файла
func (api *httpAPI) generateCryptoKeys() {
	privateKeyFile, err := os.ReadFile(config.ServerOptions.PrivateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(privateKeyFile) == 0 {
		log.Fatal(err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyFile)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	config.ServerOptions.PrivateCryptoKey = privateKey
}

// startGrpc запускает сервер grpc
func (api *httpAPI) startGrpc() {
	pb.RegisterKeeperDataServer(api.grpc, &DataServer{appRepo: api.app})
}

// closeDB вызывает функцию закрытия соединений с БД.
func (api *httpAPI) closeDB() {
	err := api.app.CloseDB()
	if err != nil {
		log.Printf("Server: error closing db, %s", err)
	}
}
