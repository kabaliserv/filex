package files

import (
	"context"
	"errors"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/files/local"
	"github.com/kabaliserv/filex/store/files/s3"
	tusd "github.com/tus/tusd/pkg/handler"
	"gorm.io/gorm"
	"io"
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
	db        *gorm.DB
	repo      *gorm.DB
	cacheRepo *gorm.DB
	options   core.StoreOption
	storage   *tusd.StoreComposer
}

func NewFileStore(db *gorm.DB, options core.StoreOption) core.FileStore {
	storage := newTusdFileStore(options)

	err := db.AutoMigrate(core.File{})
	if err != nil {
		log.Panic(err)
	}

	err = db.AutoMigrate(core.FileCache{})
	if err != nil {
		log.Panic(err)
	}

	return &fileStore{
		db:        db,
		repo:      db.Model(core.File{}),
		cacheRepo: db.Model(core.FileCache{}),
		options:   options,
		storage:   storage.GetStoreComposer(),
	}
}

func (f *fileStore) Find(filter core.File) ([]*core.File, error) {
	var users []*core.File
	result := f.repo.Where(filter).Find(&users)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return users, nil
}

func (f *fileStore) FindById(id string) (*core.File, error) {
	var file core.File
	if err := f.repo.First(&file, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

func (f *fileStore) Create(file *core.File) error {
	return f.repo.Create(file).Error
}

func (f *fileStore) CreateFromCache(id string) error {

	if f.has(id) {
		return errors.New("file is already in store")
	}

	if f.hasInCache(id) {
		return errors.New("cache file is not found")
	}

	ctx := context.Background()
	upload, err := f.storage.Core.GetUpload(ctx, id)
	if err != nil {
		return err
	}
	info, err := upload.GetInfo(ctx)
	if err != nil {
		return err
	}

	file := core.File{
		FileId:    info.ID,
		Size:      info.Size,
		Name:      info.MetaData["filename"],
		Type:      info.MetaData["filetype"],
		StorageID: info.MetaData["storage_id"],
	}

	err = f.DeleteInCacheById(id)
	if err != nil {
		return err
	}

	return f.repo.Create(&file).Error
}

func (f *fileStore) Update(file *core.File) error {
	return nil
}

func (f *fileStore) Delete(id string) error {
	ctx := context.Background()

	upload, err := f.storage.Core.GetUpload(ctx, id)
	if err != nil {
		return err
	}

	terminate := f.storage.Terminater.AsTerminatableUpload(upload)
	err = terminate.Terminate(ctx)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	err = f.repo.Where(&core.File{FileId: id}).Delete(core.File{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}

func (f *fileStore) FindInCache(filter core.FileCache) ([]*core.FileCache, error) {
	var cacheFiles []*core.FileCache
	err := f.cacheRepo.Where(filter).Find(&cacheFiles).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return cacheFiles, nil
}

func (f *fileStore) FindInCacheById(id string) (*core.FileCache, error) {
	cacheFiles, err := f.FindInCache(core.FileCache{FileId: id})
	if err != nil {
		return nil, err
	}
	return cacheFiles[0], err
}

func (f *fileStore) CreateInCache(file *core.FileCache) error {
	return f.cacheRepo.Create(file).Error
}

func (f *fileStore) UpdateInCache(file *core.FileCache) error {
	return nil
}

func (f *fileStore) DeleteInCache(file core.FileCache) error {
	return f.cacheRepo.Where(file).Delete(core.FileCache{}).Error
}

func (f *fileStore) DeleteInCacheById(fileId string) error {
	return f.DeleteInCache(core.FileCache{FileId: fileId})
}

func (f *fileStore) GetReader(id string) (io.Reader, error) {
	ctx := context.Background()
	upload, err := f.storage.Core.GetUpload(ctx, id)
	if err != nil {
		return nil, err
	}
	return upload.GetReader(ctx)
}

func (f *fileStore) TusdStoreComposer() *tusd.StoreComposer {
	return f.storage
}

func (f *fileStore) find(filter core.File) (users []*core.File, err error) {
	result := f.repo.Where(filter).Find(&users)
	if result.Error != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return
}

func (f *fileStore) has(fileId string) bool {
	var count int64
	f.repo.Where(core.File{FileId: fileId}).Count(&count)
	return count > 0
}

func (f *fileStore) hasInCache(fileId string) bool {
	var count int64
	f.db.Where(core.FileCache{FileId: fileId}).Count(&count)
	return count > 0
}
