package main

import (
	"context"
	"github.com/drone/signal"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/server"
	"github.com/sirupsen/logrus"
)

func main() {
	config, err := config.Environ()
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}

	ctx := signal.WithContext(
		context.Background(),
	)

	println("Start app...")
	app, err := InitializeApplication(config)
	if err != nil {
		panic(err)
	}

	app.server.ListenAndServe(ctx)

}

type application struct {
	server *server.Server
}

func newApplication(
	server *server.Server,
) application {
	return application{
		server: server,
	}
}
