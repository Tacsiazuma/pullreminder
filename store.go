package main

type Store interface {
	AddRepository(repo *Repository) error
	AddCredentials(provider, token string) error
	Repositories() ([]*Repository, error)
	Credentials() (map[string]string, error)
}

type FakeStore struct {
	repositories []*Repository
	credentails  map[string]string
}

func (s *FakeStore) AddRepository(repo *Repository) error {
	for _, r := range s.repositories {
		if r.Equal(repo) {
			return ErrRepositoryDuplicate
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

func (s *FakeStore) Repositories() ([]*Repository, error) {
	return s.repositories, nil
}
