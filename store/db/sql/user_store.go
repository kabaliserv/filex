package sql

import (
	"github.com/google/uuid"
	"github.com/kabaliserv/filex/core"
	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"primaryKey;uniqueIndex"`
	Username     string `gorm:"unique"`
	PasswordHash string
	Email        string
	Role         int8
}

type UserDBStore struct {
	*gorm.DB
}

func (s UserDBStore) GetUserById(id string) (*core.User, error) {

	user := User{ID: id}

	if err := s.get(&user); err != nil {
		return nil, err
	}

	userCore := s.toUserCore(user)

	return &userCore, nil

}

func (s UserDBStore) GetUserByName(name string) (*core.User, error) {

	user := User{Username: name}

	if err := s.get(&user); err != nil {
		return nil, err
	}

	userCore := s.toUserCore(user)

	return &userCore, nil

}

func (s UserDBStore) GetUserByEmail(email string) (*core.User, error) {

	user := User{Email: email}

	if err := s.get(&user); err != nil {
		return nil, err
	}

	userCore := s.toUserCore(user)

	return &userCore, nil

}

func (s UserDBStore) InsertUser(data core.User) (*core.User, error) {

	id := uuid.New()

	for {
		if !s.Has(id.String()) {
			break
		}
		id = uuid.New()
	}

	user := s.fromUserCore(data)

	user.ID = id.String()

	if err := s.create(&user); err != nil {
		return nil, err
	}

	userCore := s.toUserCore(user)

	return &userCore, nil

}

func (s UserDBStore) Has(id string) bool {

	var c int64

	s.Where(&User{ID: id}).Count(&c)

	return c > 0

}

func (s UserDBStore) create(v *User) error {

	return s.Create(v).Error

}

func (s UserDBStore) get(v *User) error {

	return s.Where(v).First(v).Error

}

func (s UserDBStore) save(id string, v map[string]interface{}) error {

	return s.Where(&User{ID: id}).Updates(v).Error

}

func (s UserDBStore) fromUserCore(u core.User) User {

	return User{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Email:        u.Email,
		Role:         u.Role,
	}

}

func (s UserDBStore) toUserCore(u User) core.User {

	return core.User{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Email:        u.Email,
		Role:         u.Role,
	}

}
