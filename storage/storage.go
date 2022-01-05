package storage

import (
	"github.com/kabaliserv/filex/cmd/kbs-filex/config"
	"github.com/kabaliserv/filex/storage/local"
	"github.com/kabaliserv/filex/storage/s3"
	tusd "github.com/tus/tusd/pkg/handler"
	"log"
	"os"
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

		path := "/data/files"
		
		f, err := os.Stat(path)

		if err != nil {

			if err := os.MkdirAll(path, os.ModeDir); err != nil {
				log.Fatalln(err)
			}

		} else if !f.IsDir() {
			log.Fatalln(path, "is not a directory")
		}

		store = local.New(path)
	}
	return store
}
