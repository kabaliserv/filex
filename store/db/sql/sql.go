package sql

import (
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/files"
	"gorm.io/gorm"
	"log"
)

type Store struct {
	db           *gorm.DB
	userStore    *UserDBStore
	accessStore  *AccessDBStore
	sessionStore *sessionStore
	fileStore    *fileStore
}

func New(option core.StoreOption) core.Store {
	db, err := getConnection(option.DatabaseDriver, option.DatabaseEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	datafileStore := files.New(option)

	return &Store{
		db:           db,
		userStore:    newUserStore(db),
		accessStore:  newAccessStore(db),
		sessionStore: newSessionStore(db, option.SessionSecret),
		fileStore:    newFileStore(db, datafileStore),
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

func (s *Store) UserStore() core.UserStore {
	return s.userStore
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

func migrateTable(db *gorm.DB, table interface{}) {
	err := db.AutoMigrate(table)
	if err != nil {
		return
	}
}
