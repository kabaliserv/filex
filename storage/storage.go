package storage

import (
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/storage/local"
	"github.com/kabaliserv/filex/storage/s3"
	tusd "github.com/tus/tusd/pkg/handler"
)

type (
	Storage interface {
		GetStoreComposer() *tusd.StoreComposer
	}
)

func New(config config.Config) Storage {
	var store Storage
	if config.Storage.S3Bucket != "" {
		store = s3.New(config.Storage.S3Bucket, config.Storage.S3EndPoint)
	} else {
		store = local.New("./tmp/files")
	}
	return store
}
