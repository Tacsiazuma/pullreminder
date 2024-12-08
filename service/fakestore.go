package service

import c "tacsiazuma/pullreminder/contract"

func NewFakeStore() Store {
	return &FakeStore{repositories: make([]*c.Repository, 0), credentails: make(map[string]string)}
}

func (s *FakeStore) AddRepository(repo *c.Repository) error {
	for _, r := range s.repositories {
		if r.Equal(repo) {
			return c.ErrRepositoryDuplicate
		}
	}
	s.repositories = append(s.repositories, repo)
	return nil
}

func (s *FakeStore) AddCredentials(provider, token string) error {
	s.credentails[provider] = token
	return nil
}

func (s *FakeStore) Credentials() (map[string]string, error) {
	return s.credentails, nil
}

func (s *FakeStore) Repositories() ([]*c.Repository, error) {
	return s.repositories, nil
}
