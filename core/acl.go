package core

import "time"

type (
	UserAccess struct {
		Token     string
		Duration  time.Duration
		CreatedAt time.Time
		UpdatedAt time.Time
		IsRevoked bool
		Save      func() error
	}

	AccessUpload struct {
		ID                int64
		AccessId          string
		MaxSize           int64
		StorageId         string
		CreatedAt         time.Time
		DeletedAt         time.Time
		AllowDeferredSize bool
		IsRevoked         bool
	}

	AccessUploadStore interface {
		Find(where AccessUpload) (accesses []*AccessUpload, err error)
		FindById(accessId string) (access *AccessUpload, err error)
		Save(access *AccessUpload) (err error)
		Delete(access *AccessUpload) (err error)
	}
)
