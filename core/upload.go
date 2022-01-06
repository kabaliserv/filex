package core

import "time"

type Upload struct {
	ID                 string
	File               File
	PasswordHash       string
	CreatedAt          time.Time
	Duration           time.Duration
	RemoveFileOnExpire bool
	Delete             func() error
}

func (u *Upload) RequireAuth() bool {
	return u.PasswordHash != ""
}
