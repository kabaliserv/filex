package users

import (
	"github.com/kabaliserv/filex/core"
	"gorm.io/gorm"
)

type userStore struct {
	db      *gorm.DB
	options core.StoreOption
}

func NewUserStore(db *gorm.DB, options core.StoreOption) core.UserStore {
	return &userStore{db: db, options: options}
}

func (u *userStore) table() *gorm.DB {
	return u.db.Table("users")
}

func (u *userStore) Find(user *core.User) ([]*core.User, error) {
	return nil, nil
}

func (u *userStore) FindOne(userId string) (*core.User, error) {
	return nil, nil
}

func (u *userStore) FindByName(username string) (*core.User, error) {
	return nil, nil
}

func (u *userStore) FindByEmail(email string) (*core.User, error) {
	return nil, nil
}

func (u *userStore) Create(user *core.User) error {
	return nil
}

func (u *userStore) Save(user *core.User) error {
	return nil
}

func (u *userStore) Has(userId string) bool {
	return false
}
