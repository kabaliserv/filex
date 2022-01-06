//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
)

func InitializeApplication(config config.Config) (application, error) {
	wire.Build(
		storeSet,
		serverSet,
		newApplication,
	)
	return application{}, nil
}
