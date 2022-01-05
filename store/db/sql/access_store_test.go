package sql

import (
	"errors"
	"github.com/kabaliserv/filex/core"
	"testing"
	"time"
)

func TestNewUserAccess(t *testing.T) {
	accessStore := New("sqlite", "/tmp/database_test.sqlite").AccessStore()
	newUserAccess, err := accessStore.NewUserAccess()

	if err == nil {
		newUserAccess.Duration = time.Hour * 4

		if newUserAccess.Save != nil {
			err = newUserAccess.Save()
		} else {
			err = errors.New("save Function is not exist")
		}
		if err == nil {

			var result interface{}
			result, err = accessStore.GetAccess(newUserAccess.Token)

			_, ok := result.(core.UserAccess)

			if !ok {
				err = errors.New("insert user access not work ")
			}

		}
	}

	if err != nil {
		t.Errorf(`Error: %v, want "", error`, err)
	}

}
