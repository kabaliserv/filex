package main

import (
	"github.com/google/wire"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/db/sql"
	"log"
)

var storeSet = wire.NewSet(
	provideStore,
	provideStoreOption,
)

func provideStore(option core.StoreOption) core.Store {
	var store core.Store

	switch option.DatabaseDriver {
	case "sqlite3", "mysql", "postgres":
		store = sql.New(option)
	default:
		log.Fatalln("invalid database driver :", option.DatabaseDriver)
	}

	return store
}

func provideStoreOption(config config.Config) core.StoreOption {
	return core.StoreOption{
		DatabaseDriver:      config.Database.Driver,
		DatabaseEndpoint:    config.Database.DataSource,
		FileStoreS3Bucket:   config.S3.Bucket,
		FileStoreS3Endpoint: config.S3.Endpoint,
		FileStoreLocalPath:  config.Storage.LocalPath,
		SessionSecret:       config.Session.Secret,
	}
}
