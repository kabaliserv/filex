package core

type (
	User struct {
		ID           string       `json:"id"`
		Username     string       `json:"username"`
		PasswordHash string       `json:"-"`
		Email        string       `json:"email"`
		Admin        bool         `json:"admin"`
		Save         func() error `json:"-"`
	}

	UserStore interface {
		Get(userId string) (*User, error)
		GetByName(username string) (*User, error)
		GetByEmail(email string) (*User, error)
		Add(user *User) error
		Has(userId string) bool
	}
)
