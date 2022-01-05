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

	switch config.Database.Driver {
	case "sqlite3", "mysql", "postgres":
		store = sql.New(config.Database.Driver, config.Database.DataSource)
	}

	return store
}
