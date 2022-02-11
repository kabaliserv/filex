package core

import (
	"github.com/google/uuid"
	tusd "github.com/tus/tusd/pkg/handler"
	"gorm.io/gorm"
	"io"
	"time"
)

type (
	File struct {
		ID        uint      `json:"-" gorm:"primaryKey"`
		UUID      string    `json:"id" gorm:"uniqueIndex"`
		FileId    string    `json:"-"`
		Name      string    `json:"name"`
		Type      string    `json:"type"`
		Size      int64     `json:"size"`
		StorageID uint      `json:"-"` // null if is uploaded by guest user
		Storage   Storage   `json:"-"`
		IsFinal   bool      `json:"-"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"-"`
		DeletedAt time.Time `json:"deleted_at"`
	}

	FileStoreComposer struct {
		*tusd.StoreComposer
	}

	FileStore interface {
		Find(where File) (files []File, err error)
		FindByID(id uint) (file File, err error)
		FindByUUID(uuid string) (file File, err error)
		GetReader(file *File) (io.Reader, error)
		Create(file *File) (err error)
		Save(file *File) (err error)
		Delete(file *File) (err error)
	}
)

func (s *File) BeforeCreate(tx *gorm.DB) (err error) {

	has := func(id string) bool {
		var count int64
		tx.Where(File{UUID: id}).Count(&count)
		return count > 0
	}

	var id string

	for {
		id = uuid.New().String()
		if !has(id) {
			break
		}
	}

	s.UUID = id
	return
}
