package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kabaliserv/filex/core"
	sessionStore "github.com/kabaliserv/filex/store/sessions"
	userStore "github.com/kabaliserv/filex/store/users"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLoginWithUsername(t *testing.T) {
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
	sessions := sessionStore.NewSessionStore(db, "Secret-123")

	p, err := bcrypt.GenerateFromPassword([]byte("C0mpleX_P@ssw0rd"), PasswordCost)

	err = users.Create(&core.User{
		Username:     "test",
		Email:        "test@gmail.com",
		PasswordHash: string(p),
	})

	if err != nil {
		t.Error(err)
	}

	f := HandleLogin(users, sessions)

	{ // test login with username
		body := `{"username":"test","password":"C0mpleX_P@ssw0rd"}`

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

}

func TestLoginWithEmail(t *testing.T) {
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
	sessions := sessionStore.NewSessionStore(db, "Secret-123")

	p, err := bcrypt.GenerateFromPassword([]byte("C0mpleX_P@ssw0rd"), PasswordCost)

	err = users.Create(&core.User{
		Username:     "test",
		Email:        "test@gmail.com",
		PasswordHash: string(p),
	})

	if err != nil {
		t.Error(err)
	}

	f := HandleLogin(users, sessions)

	{ // test login with username
		body := `{"username":"test@gmail.com","password":"C0mpleX_P@ssw0rd"}`

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
}
