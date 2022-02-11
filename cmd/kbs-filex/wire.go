//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/service/token"
)

func InitializeApplication(config config.Config) (application, error) {
	wire.Build(
		token.New,
		storeSet,
		serverSet,
		storageSet,
		newApplication,
	)
	return application{}, nil
}
