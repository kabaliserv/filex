package files

import (
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/files/local"
	"github.com/kabaliserv/filex/store/files/s3"
	tusd "github.com/tus/tusd/pkg/handler"
	"log"
	"os"
)

type (
	Storage interface {
		GetStoreComposer() *tusd.StoreComposer
	}
)

func New(option core.StoreOption) Storage {

	var store Storage

	if option.FileStoreS3Bucket != "" {

		store = s3.New(option.FileStoreS3Bucket, option.FileStoreS3Endpoint)

	} else {

		path := option.FileStoreLocalPath

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
