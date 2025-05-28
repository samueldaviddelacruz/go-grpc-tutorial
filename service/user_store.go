package service

import "sync"

type UserStore interface {
	Save(user *User) error
	Find(username string) (*User, error)
}

type InMemoryUserStore struct {
	mutex sync.Mutex
	users map[string]*User
}

func NewInMemorUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

func (store *InMemoryUserStore) Save(user *User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if store.users[user.Username] != nil {
		return ErrAlreadyExists
	}
	store.users[user.Username] = user.Clone()
	return nil
}

func (store *InMemoryUserStore) Find(username string) (*User, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	user := store.users[username]
	if user == nil {
		return nil, nil
	}

	return user.Clone(), nil
}
