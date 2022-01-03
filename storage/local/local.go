package local

import (
	"github.com/tus/tusd/pkg/filelocker"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
)

type LocalStore struct {
	Dir      string
	store    filestore.FileStore
	locker   filelocker.FileLocker
	composer *tusd.StoreComposer
}

func New(Dir string) *LocalStore {
	store := &LocalStore{
		Dir: Dir,
	}

	return store.init()
}

func (s *LocalStore) init() *LocalStore {
	s.store = filestore.New(s.Dir)
	s.locker = filelocker.New(s.Dir)

	return s
}

func (s *LocalStore) GetStoreComposer() *tusd.StoreComposer {
	if s.composer == nil {
		s.composer = tusd.NewStoreComposer()
		s.store.UseIn(s.composer)
		s.locker.UseIn(s.composer)
	}

	return s.composer
}
