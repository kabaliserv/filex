package main

import (
	"errors"
	"github.com/google/wire"
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/files"
	"github.com/kabaliserv/filex/store/sessions"
	"github.com/kabaliserv/filex/store/users"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var storeSet = wire.NewSet(
	provideStoreOption,
	provideDatabase,
	files.NewFileStore,
	users.NewUserStore,
	sessions.NewSessionStore,
)

func provideDatabase(option core.StoreOption) (*gorm.DB, error) {
	var connector gorm.Dialector

	switch option.DatabaseDriver {
	case "postgres":
		connector = postgres.Open(option.DatabaseEndpoint)
	case "mysql":
		connector = mysql.Open(option.DatabaseEndpoint)
	case "sqlite3":
		connector = sqlite.Open(option.DatabaseEndpoint)
	default:
		return nil, errors.New("invalid database driver")
	}

	return gorm.Open(connector, &gorm.Config{})
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
