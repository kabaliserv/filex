package sql

import (
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/kabaliserv/filex/core"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type AccessType string

const (
	TypeUserAccess    AccessType = "user"
	TypeStorageAccess            = "storage"
	TypeFileAccess               = "file"
)

type Access struct {
	gorm.Model
	Token     string `gorm:"uniqueIndex"`
	EntityID  string
	Type      AccessType
	Duration  time.Duration
	CreatedAt time.Time
	UpdatedAt time.Time
	IsRevoked bool
}

type UserAccess struct {
}

type AccessDBStore struct {
	*gorm.DB
}

func (s *AccessDBStore) GetAccess(token string) (interface{}, error) {
	access, err := s.getAccessByToken(token)
	if err != nil {
		return nil, err
	}
	return s.toUserAccess(access), nil
}

func (s *AccessDBStore) GetUserAccess(token string) (*core.UserAccess, error) {
	userAccess := Access{Token: token, Type: TypeUserAccess}

	if result := s.Where(&userAccess).First(&userAccess); result.Error != nil {
		return nil, result.Error
	}

	t := s.toUserAccess(userAccess)

	t.Save = s.getSaveHandle(&t)

	return &t, nil
}

func (s AccessDBStore) NewUserAccess() (*core.UserAccess, error) {

	token := uuid.New()

	for {
		ok := s.isExistAccess(token.String())

		if !ok {
			break
		}

		token = uuid.New()
	}

	access := Access{Token: token.String(), Type: TypeUserAccess}

	if err := s.insertOne(&access); err != nil {
		return &core.UserAccess{}, err
	}

	userAccess := s.toUserAccess(access)

	userAccess.Save = s.getSaveHandle(&userAccess)

	return &userAccess, nil
}

func (s AccessDBStore) getSaveHandle(access interface{}) func() error {

	f := func(a interface{}) map[string]interface{} {
		var result map[string]interface{}

		if userAccess, ok := access.(*core.UserAccess); ok {
			result = structs.Map(s.fromUserAccess(*userAccess))
		}

		return result
	}

	oldValue := f(access)

	token := oldValue["Token"].(string)

	ff := func() error {
		newValue := f(access)

		var updateValue = make(map[string]interface{})

		for k, v := range newValue {
			if !reflect.DeepEqual(oldValue[k], v) {
				updateValue[k] = v
			}
		}

		return s.updateOne(token, updateValue)
	}

	return ff
}

func (s AccessDBStore) getAccessByToken(token string) (Access, error) {
	access := Access{Token: token}

	if result := s.Where(&access).First(&access); result.Error != nil {
		return Access{}, result.Error
	}

	return access, nil
}

func (s AccessDBStore) insertOne(item *Access) error {
	return s.Create(item).Error
}

func (s AccessDBStore) updateOne(token string, data map[string]interface{}) error {
	return s.Where(&Access{Token: token}).Updates(data).Error
}

func (s AccessDBStore) isExistAccess(token string) bool {
	var count int64

	s.Where(&Access{Token: token}).Count(&count)

	return count > 0
}

func (s AccessDBStore) toUserAccess(access Access) core.UserAccess {
	return core.UserAccess{
		Token:     access.Token,
		Duration:  access.Duration,
		CreatedAt: access.CreatedAt,
		UpdatedAt: access.UpdatedAt,
		IsRevoked: access.IsRevoked,
	}
}

func (s AccessDBStore) fromUserAccess(access core.UserAccess) Access {
	return Access{
		Token:     access.Token,
		Duration:  access.Duration,
		CreatedAt: access.CreatedAt,
		UpdatedAt: access.UpdatedAt,
		IsRevoked: access.IsRevoked,
	}
}
