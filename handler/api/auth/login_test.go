package auth

import (
	"github.com/kabaliserv/filex/core"
	"github.com/kabaliserv/filex/store/db/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	db := sql.New("sqlite", "file::memory:?cache=shared")

	userDB := db.UserStore()
	sessionDB := db.SessionStore()

	p, err := bcrypt.GenerateFromPassword([]byte("C0mpleX_P@ssw0rd"), PasswordCost)

	_, err = userDB.InsertUser(core.User{
		Username:     "test",
		PasswordHash: string(p),
	})

	if err != nil {
		t.Error(err)
	}

	f := HandleLogin(userDB, sessionDB)

	body := `{"username":"test","password":"C0mpleX_P@ssw0rd"}`

	req := httptest.NewRequest("POST", "/fake", strings.NewReader(body))
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	f(w, req)

	res := w.Result()
	defer res.Body.Close()
	defer func(db core.Store) {
		err := db.Close()
		if err != nil {
			t.Error(err)
		}
	}(db)

	if httpCode := res.StatusCode; httpCode != http.StatusNoContent {
		t.Errorf("expected http code 204 got %v", httpCode)
	}
}
