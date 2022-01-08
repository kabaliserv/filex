package core

import tusd "github.com/tus/tusd/pkg/handler"

type (
	File struct {
		ID        string // uuid
		Name      string
		Type      string
		Size      int64
		StorageID string // null if is uploaded by guest user
	}

	FileCache struct {
		ID              string // uuid
		ContextUploadID string // uuid
	}

	FileStore interface {
		Find(filter *File) ([]*File, error)
		FindById(id string) (*File, error)
		Create(*File) error
		CreateFromCache(id string) error
		Update(*File) error
		Delete(*File) error
		DeleteById(id string) error

		FindInCache(filter *FileCache) ([]*FileCache, error)
		FindInCacheById(id string) (*FileCache, error)
		CreateInCache(*FileCache) error
		UpdateInCache(*FileCache) error
		DeleteInCache(*FileCache) error
		DeleteInCacheById(fileId string) error

		TusdStoreComposer() *tusd.StoreComposer
	}
)
