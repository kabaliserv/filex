package auth

import (
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/db/sql"
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

	db := sql.New(options)

	f := HandleRegister(db.UserStore())

	body := `{"username":"test","password":"C0mpleX_P@ssw0rd","email":"test1@test.com"}`

	req := httptest.NewRequest("POST", "/fake", strings.NewReader(body))
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	f(w, req)

	res := w.Result()
	defer res.Body.Close()
	defer func(db core.Store) {
		err := db.CloseConnection()
		if err != nil {
			t.Error(err)
		}
	}(db)

	if httpCode := res.StatusCode; httpCode != http.StatusNoContent {
		t.Errorf("expected http code 204 got %v", httpCode)
	}
}
