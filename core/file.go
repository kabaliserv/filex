package core

import tusd "github.com/tus/tusd/pkg/handler"

type (
	File struct {
		ID        string
		Name      string
		Type      string
		Size      int64
		StorageID string
		Delete    func() error
	}

	FileCache struct {
		FileID   string
		ClientID string
	}

	FileStore interface {
		Get(fileId string) (*File, error)
		New(fileId string) (*File, error)
		NewWithStorageId(fileId, storageId string) (*File, error)
		GetInCache(fileId string) (*FileCache, error)
		AddInCache(fileId string, clientId string) error
		HasInCache(fileId string) bool
		DelInCache(fileId string) error
		GetInCacheByClientId(clientId string) (*FileCache, error)
		GetTusdStoreComposer() *tusd.StoreComposer
	}
)
