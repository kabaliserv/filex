package sql

import (
	"github.com/kabaliserv/filex/core"
	"github.com/wader/gormstore/v2"
	"gorm.io/gorm"
	"log"
)

type Store struct {
	db          *gorm.DB
	userStore   core.UserStore
	accessStore core.AccessStore
	session     core.SessionStore
}

func New(database, endpoint string) core.Store {
	db, err := getConnection(database, endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	return &Store{db: db}
}

func (s *Store) GetDB() *gorm.DB {
	return s.db
}

func (s *Store) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (s *Store) UserStore() core.UserStore {
	if s.userStore == nil {
		s.userStore = &UserDBStore{s.db.Model(&User{})}
		migrateTable(s.db, &User{})
	}
	return s.userStore
}

func (s *Store) AccessStore() core.AccessStore {
	if s.accessStore == nil {
		s.accessStore = &AccessDBStore{s.db.Model(&Access{})}
		migrateTable(s.db, &Access{})
	}
	return s.accessStore
}

func (s *Store) SessionStore() core.SessionStore {
	if s.session == nil {
		gstore := gormstore.New(s.GetDB(), []byte("azertFRDSFV563yuiopqsdfVD36514ghjklmwxRcvbnd6f8GRSGs13rg-51h-8hhFs"))
		s.session = &SessionStore{gstore}
	}
	return s.session
}

func migrateTable(db *gorm.DB, table interface{}) {
	err := db.AutoMigrate(table)
	if err != nil {
		return
	}
}
