package upload

import (
	"errors"
	"github.com/kabaliserv/filex/core"
	gonanoid "github.com/matoous/go-nanoid"
	"gorm.io/gorm"
)

type uploadStore struct {
	db *gorm.DB
}

func NewUploadStore(db *gorm.DB, options core.StoreOption) core.UploadStore {
	return &uploadStore{db: db}
}

func (u uploadStore) Find(where core.Upload) (uploads []*core.Upload, err error) {
	result := u.db.Model(core.Upload{}).Where(where).Find(&uploads)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return
}

func (u uploadStore) FindById(id string) (upload *core.Upload, err error) {
	result := u.db.Model(core.Upload{}).First(&upload, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return
}

func (u uploadStore) Create(upload *core.Upload) (err error) {
	var id string
	for {
		id = gonanoid.MustID(32)
		if !u.has(id) {
			break
		}
	}
	upload.ID = id
	err = u.db.Model(core.Upload{}).Create(&upload).Error
	return
}

func (u uploadStore) Save(upload *core.Upload) error {
	result := u.db.Model(core.Upload{}).Where(core.Upload{ID: upload.ID}).Save(upload)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u uploadStore) Delete(upload *core.Upload) error {
	result := u.db.Model(core.Upload{}).Where(core.Upload{ID: upload.ID}).Delete(core.Upload{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u uploadStore) DeleteById(id string) error {
	return u.Delete(&core.Upload{ID: id})
}

func (u uploadStore) has(id string) bool {
	var count int64

	u.db.Model(core.Upload{}).Where(core.Upload{ID: id}).Count(&count)

	return count == 1
}
