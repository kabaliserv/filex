package storage

import (
	"errors"
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type storageStore struct {
	db *gorm.DB
}

func NewStorageStore(db *gorm.DB, options core.StoreOption) core.StorageStore {
	if err := db.AutoMigrate(&core.Storage{}); err != nil {
		log.Panicf("error on migrate users table: %#v", err)
	}
	return &storageStore{db: db}
}

func (s storageStore) Find(where core.Storage) (storages []core.Storage, err error) {
	err = s.db.Where(where).Find(&storages).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = core.ErrNotFound
	}

	return
}

func (s storageStore) FindByID(id uint) (storage core.Storage, err error) {
	err = s.db.Model(storage).First(&storage, id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = core.ErrNotFound
	}

	return
}

func (s storageStore) FindByUUID(uuid string) (storage core.Storage, err error) {
	err = s.db.Model(storage).Where(core.Storage{UUID: uuid}).Take(&storage).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = core.ErrNotFound
	}

	return
}

func (s storageStore) Save(storage *core.Storage) (err error) {
	return s.db.Save(storage).Error
}

func (s storageStore) Delete(storage *core.Storage) (err error) {
	return s.db.Where(storage).Delete(storage).Error
}
