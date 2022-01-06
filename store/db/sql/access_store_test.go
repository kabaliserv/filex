package sql

import (
	"errors"
	"github.com/kabaliserv/filex/core"
	"testing"
	"time"
)

func TestNewUserAccess(t *testing.T) {
	options := core.StoreOption{
		FileStoreLocalPath: "/tmp/files",
		DatabaseDriver:     "sqlite3",
		DatabaseEndpoint:   "file::memory:?cache=shared",
	}
	db := New(options)
	accessStore := db.AccessStore()
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

	defer func(db core.Store) {
		err := db.CloseConnection()
		if err != nil {
			t.Error(err)
		}
	}(db)

	if err != nil {
		t.Errorf(`Error: %v, want "", error`, err)
	}

}
