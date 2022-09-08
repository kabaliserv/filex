package access

import (
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type accessUploadStore struct {
	db *gorm.DB
}

func NewAccessUploadStore(db *gorm.DB, options core.StoreOption) core.AccessUploadStore {
	if err := db.AutoMigrate(core.AccessUpload{}); err != nil {
		log.Panicf("autoMigrate(access upload): %s", err)
	}
	return &accessUploadStore{db: db}
}

func (a accessUploadStore) Find(where core.AccessUpload) (accesses []*core.AccessUpload, err error) {
	err = a.db.Where(where).Find(accesses).Error
	return
}

func (a accessUploadStore) FindById(accessId string) (access *core.AccessUpload, err error) {
	err = a.db.Where(core.AccessUpload{AccessId: accessId}).Find(access).Error
	return
}

func (a accessUploadStore) Save(access *core.AccessUpload) (err error) {
	err = a.db.Where(core.AccessUpload{ID: access.ID}).Save(access).Error
	return
}

func (a accessUploadStore) Delete(access *core.AccessUpload) (err error) {
	err = a.db.Where(core.AccessUpload{ID: access.ID}).Delete(access).Error
	return
}
