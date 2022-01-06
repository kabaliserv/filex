package sql

import (
	"context"
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/files"
	"github.com/prometheus/common/log"
	tusd "github.com/tus/tusd/pkg/handler"
	"gorm.io/gorm"
)

type fileSchema struct {
	ID        string `gorm:"primarykey"`
	StorageID string
	Name      string
	Type      string
	Size      int64
}

type cacheSchema struct {
	gorm.Model
	FileId string
}

type fileStore struct {
	*gorm.DB
	storage *tusd.StoreComposer
	cache   *gorm.DB
}

func newFileStore(db *gorm.DB, storage files.Storage) *fileStore {
	table := db.Table("files")
	table.AutoMigrate(&fileSchema{})

	tableCache := db.Table("cache_files")
	tableCache.AutoMigrate(&cacheSchema{})

	return &fileStore{
		DB:      table,
		cache:   tableCache,
		storage: storage.GetStoreComposer(),
	}
}

func (f *fileStore) Get(id string) (*core.File, error) {
	return nil, nil
}

func (f *fileStore) New(fileId string) (*core.File, error) {
	return f.newFile(fileId, "")
}

func (f *fileStore) NewWithStorageId(fileId, storageId string) (*core.File, error) {
	return f.newFile(fileId, storageId)
}

func (f *fileStore) AddInCache(fileId string) error {
	return f.cache.Create(&cacheSchema{FileId: fileId}).Error
}

func (f *fileStore) GetTusdStoreComposer() *tusd.StoreComposer {
	return f.storage
}

func (f *fileStore) newFile(fileId string, storageId string) (*core.File, error) {
	info, err := f.getTusFileInfo(fileId)
	if err != nil {
		return nil, err
	}
	lFile := fileSchema{
		ID:        fileId,
		Size:      info.Size,
		Name:      info.MetaData["filename"],
		Type:      info.MetaData["filetype"],
		StorageID: storageId,
	}

	if err := f.DB.Create(&lFile).Error; err != nil {
		return nil, err
	}

	file := core.File{}

	f.fileSchemaToCoreFile(lFile, &file)

	return &file, nil
}

func (f *fileStore) getTusFileInfo(id string) (tusd.FileInfo, error) {
	ctx := context.Background()
	upload, err := f.storage.Core.GetUpload(ctx, id)
	if err != nil {
		return tusd.FileInfo{}, err
	}
	return upload.GetInfo(ctx)
}

func (f *fileStore) fileSchemaToCoreFile(from fileSchema, to *core.File) {
	to.ID = from.ID
	to.Name = from.Name
	to.Size = from.Size
	to.Type = from.Type
	to.StorageID = from.StorageID
}

func (f *fileStore) fileIsInCache(id string) bool {
	return has(f.cache, &cacheSchema{FileId: id})
}

func (f *fileStore) hasFile(id string) bool {
	return has(f.DB, &fileSchema{ID: id})
}

func create(db *gorm.DB, value interface{}) error {
	return db.Create(value).Error
}

func update(db *gorm.DB, where interface{}, value map[string]interface{}) error {
	return db.Where(where).Updates(value).Error
}

func has(db *gorm.DB, where interface{}) bool {
	var count int64

	if result := db.Where(where).Count(&count); result.Error != nil {
		log.Error(result.Error)
	}

	return count == 1
}
