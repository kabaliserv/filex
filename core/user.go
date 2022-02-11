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
	errEmailLen     = errors.New("Invalid email length")
	errEmailChar    = errors.New("Invalid character in email")
)

type (
	User struct {
		ID           uint    `json:"-" gorm:"primaryKey"`
		UUID         string  `json:"id" gorm:"uniqueIndex"`
		Avatar       string  `json:"avatar"`
		Login        string  `json:"login"`
		PasswordHash string  `json:"-"`
		Email        string  `json:"email" group:"personal"`
		Active       bool    `json:"active" group:"admin"`
		Admin        bool    `json:"admin"`
		Storage      Storage `json:"storage"`
	}

	UserStore interface {
		Find(filter User) (users []User, err error)
		FindByID(id uint) (user User, err error)
		FindByUUID(uuid string) (user User, err error)
		FindByLogin(username string) (user User, err error)
		FindByEmail(email string) (user User, err error)
		Create(user *User) error
		Save(user *User) error
		Has(userId uint) bool
	}
)

func (u *User) Validate() error {
	switch {
	case !govalidator.IsByteLength(u.Login, 1, 50):
		return errUsernameLen
	case !govalidator.Matches(u.Login, "^[.a-zA-Z0-9_-]+$"):
		return errUsernameChar
	case !govalidator.IsByteLength(u.Email, 5, 100):
		return errEmailLen
	case !govalidator.Matches(u.Email, "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"):
		return errEmailChar
	default:
		return nil
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Active = true
	err = u.Validate()

	has := func(id string) bool {
		var count int64
		tx.Where(User{UUID: id}).Count(&count)
		return count > 0
	}

	var id string

	for {
		id = uuid.New().String()
		if !has(id) {
			break
		}
	}

	u.UUID = id
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	counts := func() int64 {
		var count int64
		tx.Model(&User{}).Count(&count)
		return count
	}
	if counts() == 1 {
		tx.Model(u).Update("admin", true)
	}
	return
}
