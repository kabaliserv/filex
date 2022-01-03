//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/storage"
)

func InitializeApplication(config config.Config) (application, error) {
	wire.Build(
		storage.New,
		serverSet,
		newApplication,
	)
	return application{}, nil
}
