package core

import (
	tusd "github.com/tus/tusd/pkg/handler"
	"time"
)

type (
	File struct {
		ID        int       `json:"-"`
		FileId    string    `json:"id"`
		Name      string    `json:"name"`
		Type      string    `json:"type"`
		Size      int64     `json:"size"`
		StorageID string    `json:"-"` // null if is uploaded by guest user
		CreatedAt time.Time `json:"-"`
		UpdatedAt time.Time `json:"-"`
		DeletedAt time.Time `json:"-"`
	}

	FileCache struct {
		ID              int
		FileId          string
		ContextUploadID string // uuid
	}

	FileStore interface {
		Find(filter File) ([]*File, error)
		FindById(id string) (*File, error)
		Create(*File) error
		CreateFromCache(id string) error
		Update(*File) error
		Delete(id string) error

		FindInCache(filter FileCache) ([]*FileCache, error)
		FindInCacheById(id string) (*FileCache, error)
		CreateInCache(*FileCache) error
		UpdateInCache(*FileCache) error
		DeleteInCache(FileCache) error
		DeleteInCacheById(fileId string) error

		TusdStoreComposer() *tusd.StoreComposer
	}
)
