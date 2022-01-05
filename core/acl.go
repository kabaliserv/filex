package core

import "time"

type (
	UserAccess struct {
		Token     string
		Duration  time.Duration
		CreatedAt time.Time
		UpdatedAt time.Time
		IsRevoked bool
		Save      func() error
	}
)
