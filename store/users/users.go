package users

import (
	"errors"
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
	return &userStore{db: db, options: options}
}

func (u *userStore) Find(filter core.User) (users []core.User, err error) {
	result := u.db.Model(users).Where(filter).Preload("Storage").Find(&users)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return users, nil
}

func (u *userStore) FindByUUID(uuid string) (user core.User, err error) {

	err = u.db.Model(user).Where(core.User{UUID: uuid}).Preload("Storage").Find(&user).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = core.ErrNotFound
	}

	return
}

func (u *userStore) FindByID(id uint) (user core.User, err error) {

	err = u.db.Model(user).Preload("Storage").First(&user, id).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = core.ErrNotFound
	}

	return
}

func (u *userStore) FindByLogin(login string) (user core.User, err error) {

	err = u.db.Model(user).Where(core.User{Login: login}).Preload("Storage").Find(&user).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = core.ErrNotFound
	}

	return
}

func (u *userStore) FindByEmail(email string) (user core.User, err error) {

	err = u.db.Model(user).Where(core.User{Email: email}).Preload("Storage").Find(&user).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = core.ErrNotFound
	}

	return
}

func (u *userStore) Create(user *core.User) (err error) {
	err = u.db.Create(user).Error

	return
}

func (u *userStore) Save(user *core.User) (err error) {
	return u.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error
}

func (u *userStore) Has(id uint) bool {
	var count int64
	u.db.Where(core.User{ID: id}).Count(&count)
	return count > 0
}
