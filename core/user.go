package core

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	errUsernameLen  = errors.New("Invalid username length")
	errUsernameChar = errors.New("Invalid character in username")
)

type (
	User struct {
		ID           int64       `json:"-"`
		UUID         uuid.UUID   `json:"id"`
		Login        string      `json:"login"`
		PasswordHash string      `json:"-"`
		Email        string      `json:"email"`
		Active       bool        `json:"active"`
		Admin        bool        `json:"admin"`
		Storage      UserStorage `json:"storage"`
		StorageID    int64       `json:"-"`
	}

	UserStorage struct {
		ID          int64     `json:"-"`
		UUID        uuid.UUID `json:"id"`
		UserID      int64     `json:"-"`
		Size        int64     `json:"size"`
		Quota       int64     `json:"quota"`
		EnableQuota bool      `json:"enable_quota"`
	}

	UserStore interface {
		Find(filter User) ([]*User, error)
		FindByID(id interface{}) (*User, error)
		FindByUUID(uuid uuid.UUID) (*User, error)
		FindByLogin(username string) (*User, error)
		FindByEmail(email string) (*User, error)
		Create(user *User) error
		Save(user *User) error
		Has(userId string) bool
	}
)

func (u *User) Validate() error {
	switch {
	case !govalidator.IsByteLength(u.Login, 1, 50):
		return errUsernameLen
	case !govalidator.Matches(u.Login, "^[.a-zA-Z0-9_-]+$"):
		return errUsernameChar
	default:
		return nil
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()
	err = u.Validate()
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		tx.Model(u).Update("admin", true)
	}
	return
}

func (u *UserStorage) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()
	return
}
