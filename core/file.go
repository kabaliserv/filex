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

	FileStore interface {
		Get(fileId string) (*File, error)
		New(fileId string) (*File, error)
		NewWithStorageId(fileId, storageId string) (*File, error)
		AddInCache(fileId string) error
		GetTusdStoreComposer() *tusd.StoreComposer
	}
)
