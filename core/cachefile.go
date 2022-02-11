package core

import (
	"time"
)

type (
	CacheFile struct {
		ID        uint      `json:"-"`
		UUID      string    `json:"id" gorm:"uniqueIndex"`
		FileId    string    `json:"-"`
		Name      string    `json:"name"`
		Type      string    `json:"type"`
		Size      int64     `json:"size"`
		StorageID uint      `json:"-"`
		CreatedAt time.Time `json:"created_at"`
		DeletedAt time.Time `json:"deleted_at"`
	}

	CacheFileStore interface {
		Find(where CacheFile) (cacheFiles []*CacheFile, err error)
		FindByID(id uint) (cacheFile *File, err error)
		FindByUUID(uuid string) (cacheFile *File, err error)
		Create(cacheFile *CacheFile) (err error)
		Save(cacheFile *CacheFile) (err error)
		Delete(uuid string) (err error)
	}
)
