package sql

import (
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/files"
	"gorm.io/gorm"
	"log"
)

type Store struct {
	db *gorm.DB
	//userStore      *UserDB
	accessStore    *AccessDBStore
	sessionStore   *sessionStore
	fileStore      *fileStore
	storageStorage *storageDB
}

func New(option core.StoreOption) core.Store {
	db, err := getConnection(option.DatabaseDriver, option.DatabaseEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	datafileStore := files.New(option)

	return &Store{
		db: db,
		//userStore:      newUserStore(db),
		accessStore:    newAccessStore(db),
		sessionStore:   newSessionStore(db, option.SessionSecret),
		fileStore:      newFileStore(db, datafileStore),
		storageStorage: newStorageStore(db),
	}
}

func (s *Store) GetDB() *gorm.DB {
	return s.db
}

func (s *Store) CloseConnection() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (s *Store) AccessStore() core.AccessStore {
	return s.accessStore
}

func (s *Store) SessionStore() core.SessionStore {
	return s.sessionStore
}

func (s *Store) FileStore() core.FileStore {
	return s.fileStore
}

func (s *Store) StorageStore() core.StorageStore {
	return s.storageStorage
}

func migrateTable(db *gorm.DB, table interface{}) {
	err := db.AutoMigrate(table)
	if err != nil {
		return
	}
}
