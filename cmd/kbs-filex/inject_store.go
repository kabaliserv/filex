package main

import (
	"github.com/google/wire"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/db/sql"
)

var storeSet = wire.NewSet(
	provideStore,
)

func provideStore(config config.Config) core.Store {
	var store core.Store

	if config.Database.Host == "" {
		store = sql.New("sqlite", "./tmp/database.sqlite")
	}

	return store
}
