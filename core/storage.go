package core

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Storage struct {
		ID          uint   `json:"-" gorm:"primaryKey"`
		UUID        string `json:"id" gorm:"uniqueIndex"`
		UserID      uint   `json:"-"`
		Size        int64  `json:"size"`
		Quota       int64  `json:"quota"`
		QuotaEnable bool   `json:"activeQuota"`
	}

	StorageStore interface {
		Find(where Storage) (storages []Storage, err error)
		FindByID(id uint) (storage Storage, err error)
		FindByUUID(uuid string) (storage Storage, err error)
		Save(storage *Storage) (err error)
		Delete(storage *Storage) (err error)
	}
)

func (s *Storage) BeforeCreate(tx *gorm.DB) (err error) {

	has := func(id string) bool {
		var count int64
		tx.Where(Storage{UUID: id}).Count(&count)
		return count > 0
	}

	var id string

	for {
		id = uuid.New().String()
		if !has(id) {
			break
		}
	}

	s.UUID = id
	return
}
