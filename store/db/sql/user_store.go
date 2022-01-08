package sql

import (
	"errors"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/kabaliserv/filex/core"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userSchema struct {
	ID           string `json:"id" gorm:"primaryKey;uniqueIndex""`
	Username     string `json:"username" gorm:"unique"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
	Admin        bool   `json:"admin"`
}

type UserDB struct {
	table *gorm.DB
}

func newUserStore(db *gorm.DB) *UserDB {
	table := db.Table("users")
	table.AutoMigrate(&userSchema{})
	return &UserDB{table}
}

func (s *UserDB) Get(userId string) (*core.User, error) {
	return s.get(userSchema{ID: userId})
}

func (s *UserDB) GetByName(name string) (*core.User, error) {
	return s.get(userSchema{Username: name})
}

func (s *UserDB) GetByEmail(email string) (*core.User, error) {
	return s.get(userSchema{Email: email})
}

func (s *UserDB) Add(user *core.User) error {

	id := uuid.New()

	for {
		if !s.Has(id.String()) {
			break
		}
		id = uuid.New()
	}

	var userSH userSchema

	if err := structToStruct(*user, &userSH); err != nil {
		return err
	}

	userSH.ID = id.String()
	userSH.PasswordHash = user.PasswordHash

	if err := s.table.Create(&userSH).Error; err != nil {
		return err
	}

	if err := structToStruct(userSH, user); err != nil {
		return err
	}

	if err := s.injectSaveHandle(user); err != nil {
		return err
	}

	return nil
}

func (s *UserDB) Has(userId string) bool {
	var c int64
	s.table.Where(&userSchema{ID: userId}).Count(&c)
	return c > 0
}

func (s *UserDB) get(where userSchema) (*core.User, error) {

	var userSH userSchema
	err := s.table.Where(&where).Find(&userSH).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var user core.User

	if err := structToStruct(&userSH, &user); err != nil {
		return nil, err
	}

	if err := s.injectSaveHandle(&user); err != nil {
		return nil, err
	}

	user.PasswordHash = userSH.PasswordHash

	return &user, nil
}

func (s *UserDB) injectSaveHandle(user *core.User) error {
	id := user.ID
	getValue := func() map[string]interface{} {
		var userSH userSchema
		if err := structToStruct(user, &userSH); err != nil {
			log.Error(err)
		}

		return structs.Map(userSH)
	}

	saveValue := func(v map[string]interface{}) error {
		return s.table.Where(&userSchema{ID: id}).Save(v).Error
	}

	user.Save = getSaveChangeFunc(saveValue, getValue)

	return nil
}
