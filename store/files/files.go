package files

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/kabaliserv/filex/core"

	tusd "github.com/tus/tusd/pkg/handler"
	"gorm.io/gorm"
)

type fileStore struct {
	db      *gorm.DB
	options core.StoreOption
	storage *tusd.StoreComposer
}

func NewFileStore(db *gorm.DB, options core.StoreOption, fileStorage core.FileStoreComposer) core.FileStore {

	err := db.AutoMigrate(core.File{})
	if err != nil {
		log.Panic(err)
	}

	return &fileStore{
		db:      db,
		options: options,
		storage: fileStorage.StoreComposer,
	}
}

func (f *fileStore) Find(where core.File) (files []core.File, err error) {

	err = f.db.Model(files).Where(where).Find(&files).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = core.ErrNotFound
	}

	return
}

func (f *fileStore) FindByID(id uint) (file core.File, err error) {

	err = f.db.Model(file).Where(core.File{ID: id}).Find(&file).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = core.ErrNotFound
	}

	return
}

func (f *fileStore) FindByUUID(id string) (file core.File, err error) {

	err = f.db.Model(file).Where(core.File{UUID: id}).Find(&file).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = core.ErrNotFound
	}

	return
}

func (f *fileStore) GetReader(file *core.File) (io.Reader, error) {
	ctx := context.Background()
	upload, err := f.storage.Core.GetUpload(ctx, file.FileId)
	if err != nil {
		return nil, err
	}
	return upload.GetReader(ctx)
}

func (f *fileStore) Create(file *core.File) error {
	return f.db.Create(file).Error
}

func (f *fileStore) Save(file *core.File) (err error) {
	return f.db.Save(file).Error
}

func (f *fileStore) Delete(file *core.File) (err error) {
	ctx := context.Background()

	var upload tusd.Upload
	upload, err = f.storage.Core.GetUpload(ctx, file.FileId)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return core.ErrNotFound
		}
		return
	}

	terminate := f.storage.Terminater.AsTerminatableUpload(upload)

	err = terminate.Terminate(ctx)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return core.ErrNotFound
		}
		return
	}

	err = f.db.Where(file).Delete(core.File{}).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return core.ErrNotFound
	}

	return
}

func (f *fileStore) TusdStoreComposer() *tusd.StoreComposer {
	return f.storage
}
