package users

import (
	"errors"
	"github.com/google/uuid"
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userStore struct {
	db      *gorm.DB
	options core.StoreOption
}

func NewUserStore(db *gorm.DB, options core.StoreOption) core.UserStore {

	if err := db.AutoMigrate(&core.User{}); err != nil {
		log.Panicf("error on migrate users table: %#v", err)
	}
	if err := db.AutoMigrate(&core.UserStorage{}); err != nil {
		log.Panicf("error on migrate user_storages table: %#v", err)
	}
	return &userStore{db: db, options: options}
}

func (u *userStore) table() *gorm.DB {
	return u.db.Table("users").Joins("Storage")
}

func (u *userStore) Find(filter core.User) ([]*core.User, error) {
	var users []*core.User
	result := u.table().Where(filter).Find(&users)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return users, nil
}

func (u *userStore) FindByID(id interface{}) (user *core.User, err error) {

	if err := u.table().First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return
}

func (u *userStore) FindByUUID(id uuid.UUID) (*core.User, error) {
	users, err := u.Find(core.User{UUID: id})
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func (u *userStore) FindByLogin(login string) (*core.User, error) {
	users, err := u.Find(core.User{Login: login})
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func (u *userStore) FindByEmail(email string) (*core.User, error) {
	users, err := u.Find(core.User{Email: email})
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func (u *userStore) Create(user *core.User) error {
	if err := u.table().Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (u *userStore) Save(user *core.User) error {
	return nil
}

func (u *userStore) Has(userId string) bool {
	return false
}
