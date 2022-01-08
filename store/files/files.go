package files

import (
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/files/local"
	"github.com/kabaliserv/filex/store/files/s3"
	tusd "github.com/tus/tusd/pkg/handler"
	"gorm.io/gorm"
	"log"
	"os"
)

type (
	tusdFileStorage interface {
		GetStoreComposer() *tusd.StoreComposer
	}
)

func newTusdFileStore(option core.StoreOption) tusdFileStorage {

	var store tusdFileStorage

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

//###############################################################################################################

type fileStore struct {
	db      *gorm.DB
	options core.StoreOption
	files   tusdFileStorage
}

func NewFileStore(db *gorm.DB, options core.StoreOption) core.FileStore {
	return &fileStore{db: db, options: options, files: newTusdFileStore(options)}
}

func (f *fileStore) table() *gorm.DB {
	return f.db.Table("files")
}

func (f *fileStore) Find(filter *core.File) ([]*core.File, error) {
	return nil, nil
}

func (f *fileStore) FindById(id string) (*core.File, error) {
	return nil, nil
}

func (f *fileStore) Create(file *core.File) error {
	return nil
}

func (f *fileStore) CreateFromCache(id string) error {
	return nil
}

func (f *fileStore) Update(file *core.File) error {
	return nil
}

func (f *fileStore) Delete(file *core.File) error {
	return nil
}

func (f *fileStore) DeleteById(id string) error {
	return nil
}

func (f *fileStore) FindInCache(filter *core.FileCache) ([]*core.FileCache, error) {
	return nil, nil
}

func (f *fileStore) FindInCacheById(id string) (*core.FileCache, error) {
	return nil, nil
}

func (f *fileStore) CreateInCache(file *core.FileCache) error {
	return nil
}

func (f *fileStore) UpdateInCache(file *core.FileCache) error {
	return nil
}

func (f *fileStore) DeleteInCache(file *core.FileCache) error {
	return nil
}

func (f *fileStore) DeleteInCacheById(fileId string) error {
	return nil
}

func (f *fileStore) TusdStoreComposer() *tusd.StoreComposer {
	return nil
}
