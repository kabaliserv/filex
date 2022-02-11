package storage

import (
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/storage/local"
	"github.com/kabaliserv/filex/storage/s3"
	log "github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/pkg/handler"
	"os"
)

type tusdFileStorage interface {
	GetStoreComposer() *tusd.StoreComposer
}

func New(options core.StoreOption) core.FileStoreComposer {

	var store tusdFileStorage

	if options.FileStoreS3Bucket != "" {

		store = s3.New(options.FileStoreS3Bucket, options.FileStoreS3Endpoint)

	} else {

		path := options.FileStoreLocalPath

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

	return core.FileStoreComposer{StoreComposer: store.GetStoreComposer()}
}
