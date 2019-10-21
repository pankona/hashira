package store

import (
	"fmt"

	"github.com/pankona/hashira/auth/user"
)

type MemStore struct {
	users map[string]*user.User
}

func NewMemStore() *MemStore {
	return &MemStore{
		users: make(map[string]*user.User),
	}
}

func (m *MemStore) Store(u *user.User) error {
	if u.ID == "" {
		return fmt.Errorf("failed to store user. user.ID must be specified. user.ID: %v", u.ID)
	}
	m.users[u.ID] = u
	return nil
}

func (m *MemStore) Fetch(userID string) (*user.User, error) {
	return m.users[userID], nil
}

func (m *MemStore) FetchByAccessToken(accesstoken string) (*user.User, error) {
	for _, v := range m.users {
		if v.AccessToken == accesstoken {
			return v, nil
		}
	}
	return nil, nil

}

func (m *MemStore) FetchByTwitterIDToken(idtoken string) (*user.User, error) {
	for _, v := range m.users {
		if v.TwitterIDToken == idtoken {
			return v, nil
		}
	}
	return nil, nil
}

func (m *MemStore) FetchByGoogleIDToken(idtoken string) (*user.User, error) {
	for _, v := range m.users {
		if v.GoogleIDToken == idtoken {
			return v, nil
		}
	}
	return nil, nil
}
