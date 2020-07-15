package storage

type Storage struct {
	Token TokenStorage
}

func NewStorage() *Storage {
	return &Storage{
		Token: newTokenStorage(),
	}
}

func newTokenStorage() TokenStorage {
	return NewTokenStorage()
}
