package core

type (
	User struct {
		ID           string
		Username     string
		PasswordHash string
		Email        string
		Role         int8
	}
)
