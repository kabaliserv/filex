package auth

import (
	"github.com/kabaliserv/filex/core"
	storageStore "github.com/kabaliserv/filex/store/storage"
	userStore "github.com/kabaliserv/filex/store/users"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGoodRequest(t *testing.T) {
	options := core.StoreOption{
		FileStoreLocalPath: "/tmp/files",
		DatabaseDriver:     "sqlite3",
		DatabaseEndpoint:   "file::memory:?cache=shared",
	}

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Error(err)
		return
	}

	rawDB, err := db.DB()
	if err != nil {
		t.Error(err)
		return
	}
	defer rawDB.Close()

	users := userStore.NewUserStore(db, options)
	_ = storageStore.NewStorageStore(db, options)

	f := HandleRegister(users)

	body := `{"login":"test","password":"C0mpleX_P@ssw0rd","email":"test1@test.com"}`

	req := httptest.NewRequest("POST", "/fake", strings.NewReader(body))
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	f(w, req)

	res := w.Result()
	defer res.Body.Close()

	if httpCode := res.StatusCode; httpCode != http.StatusNoContent {
		t.Errorf("expected http code 204 got %v", httpCode)
	}
}
