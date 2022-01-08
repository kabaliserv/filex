package core

type (
	User struct {
		ID           string      `json:"id"`
		Username     string      `json:"username"`
		PasswordHash string      `json:"-"`
		Email        string      `json:"email"`
		Admin        bool        `json:"admin"`
		Storage      UserStorage `json:"storage"`
	}

	UserStorage struct {
		ID          string `json:"id"`
		UserID      string `json:"user_id"`
		Size        int64  `json:"size"`
		Quota       int64  `json:"quota"`
		EnableQuota bool   `json:"enable_quota"`
	}

	UserStore interface {
		Find(*User) ([]*User, error)
		FindOne(userId string) (*User, error)
		FindByName(username string) (*User, error)
		FindByEmail(email string) (*User, error)
		Create(user *User) error
		Save(user *User) error
		Has(userId string) bool
	}
)
