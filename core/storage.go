package core

type (
	Storage struct {
		ID     string       `json:"id"`
		UserID string       `json:"user_id"`
		Size   int64        `json:"size"`
		Quota  int64        `json:"quota"`
		Save   func() error `json:"-"`
	}

	StorageStore interface {
		Get(storageId string) (*Storage, error)
		Add(storage *Storage) error
		Del(storageId string) error
		Has(storageId string) bool
		GetByUserId(userId string) (*Storage, error)
	}
)
