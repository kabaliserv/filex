package core

import "time"

type Upload struct {
	ID                 string
	File               File
	PasswordHash       string
	CreatedAt          time.Time
	Duration           time.Duration
	Active             bool
	RemoveFileOnExpire bool
}

type UploadStore interface {
	Find(where Upload) ([]*Upload, error)
	FindById(id string) (*Upload, error)
	Create(upload *Upload) error
	Save(upload *Upload) error
	Delete(upload *Upload) error
	DeleteById(id string) error
}

func (u *Upload) RequireAuth() bool {
	return u.PasswordHash != ""
}

type UploadOption struct {
	GuestAllow         bool
	GuestMaxUploadSize int64
}
