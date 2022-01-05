package main

import (
	"context"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/server"
	"github.com/kabaliserv/filex/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	println("Hello World!")
	config, err := config.Environ()
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}

	app, err := InitializeApplication(config)

	if err != nil {
		panic(err)
	}
	app.server.ListenAndServe(context.Background())

}

type application struct {
	server  *server.Server
	storage storage.Storage
}

func newApplication(
	server *server.Server,
	storage storage.Storage,
) application {
	return application{
		storage: storage,
		server:  server,
	}
}
