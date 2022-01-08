package sql

import (
	"errors"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type storageSchema struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Size      int64     `json:"size"`
	Quota     int64     `json:"quota"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type storageDB struct {
	table *gorm.DB
}

func newStorageStore(db *gorm.DB) *storageDB {
	table := db.Table("storage")
	table.AutoMigrate(&storageSchema{})

	return &storageDB{table}
}

func (s *storageDB) Get(storageId string) (*core.Storage, error) {
	return s.get(storageSchema{ID: storageId})
}

func (s *storageDB) GetByUserId(userId string) (*core.Storage, error) {
	return s.get(storageSchema{UserId: userId})
}

func (s *storageDB) Add(storage *core.Storage) error {
	id := uuid.New()

	for {
		if !s.Has(id.String()) {
			break
		}
		id = uuid.New()
	}

	storageSH := storageSchema{}

	if err := structToStruct(storage, &storageSH); err != nil {
		return err
	}

	storageSH.ID = id.String()

	if err := s.table.Create(&storageSH).Error; err != nil {
		return err
	}

	if err := structToStruct(&storageSH, storage); err != nil {
		return err
	}

	if err := s.injectSaveHandle(storage); err != nil {
		return err
	}

	return nil
}

func (s *storageDB) Del(storageId string) error {
	return s.table.Delete(&storageSchema{ID: storageId}).Error
}

func (s *storageDB) Has(storageId string) bool {
	var count int64
	if err := s.table.Where(&storageSchema{ID: storageId}).Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		log.Error(err)
	}
	return count == 1
}

func (s *storageDB) get(where storageSchema) (*core.Storage, error) {

	var storageSH storageSchema
	if err := s.table.Where(&where).First(&storageSH).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	storage := core.Storage{}

	if err := structToStruct(&storageSH, &storage); err != nil {
		return nil, err
	}

	s.injectSaveHandle(&storage)
	return &storage, nil
}

func (s *storageDB) injectSaveHandle(storage *core.Storage) error {
	id := storage.ID
	getValue := func() map[string]interface{} {
		s := storageSchema{}
		if err := structToStruct(storage, &s); err != nil {
			log.Error(err)
		}

		return structs.Map(s)
	}

	saveValue := func(v map[string]interface{}) error {
		return s.table.Where(&storageSchema{ID: id}).Save(v).Error
	}

	storage.Save = getSaveChangeFunc(saveValue, getValue)

	return nil
}
