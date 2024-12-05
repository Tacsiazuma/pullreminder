package main

import (
	"context"
	"errors"
)

type Service struct {
	store    Store
	provider Provider
}

func New(provider Provider) *Service {
	return &Service{
		store:    NewFakeStore(),
		provider: provider,
	}
}

var (
	ErrNoRepositoriesProvided         = errors.New("No repositories has been provided")
	ErrRepositoryMissingName          = errors.New("Repository should have a name")
	ErrRepositoryMissingOwner         = errors.New("Repository should have an owner")
	ErrRepositoryMissingProvider      = errors.New("Repository should have an provider")
	ErrRepositoryInvalidProvider      = errors.New("Providers should be (github or gitlab)")
	ErrInvalidProvider                = errors.New("Providers should be (github or gitlab)")
	ErrRepositoryDuplicate            = errors.New("Cannot add the same repository twice")
	ErrNoCredentialsProvidedForGithub = errors.New("No credentials provided for github")
	ErrNoCredentialsProvidedForGitlab = errors.New("No credentials provided for gitlab")
	ErrCannotQueryRepository          = errors.New("Repository cannot be queried")
)

func (s *Service) NeedsAttention(ctx context.Context) ([]*Pullrequest, error) {
	repos, _ := s.store.Repositories()
	if len(repos) == 0 {
		return nil, ErrNoRepositoriesProvided
	}
	creds, _ := s.store.Credentials()
	total := make([]*Pullrequest, 0)
	for _, repo := range repos {
		if creds[repo.Provider] == "" {
			switch repo.Provider {
			case "gitlab":
				return nil, ErrNoCredentialsProvidedForGitlab
			case "github":
				return nil, ErrNoCredentialsProvidedForGithub
			}
		}
		prs, err := s.provider.GetPullRequests(ctx, *repo, creds[repo.Provider], "main")
		if err != nil {
			return total, err
		}
		total = append(total, prs...)
	}
	return total, nil
}

func (s *Service) Repositories() ([]*Repository, error) {
	return s.store.Repositories()
}

func (s *Service) AddCredentials(provider, token string) error {
	switch provider {
	case "github":
		break
	case "gitlab":
		break
	default:
		return ErrInvalidProvider
	}
	return s.store.AddCredentials(provider, token)
}

func (s *Service) AddRepository(repo *Repository) error {
	if repo.Name == "" {
		return ErrRepositoryMissingName
	}
	if repo.Owner == "" {
		return ErrRepositoryMissingOwner
	}
	switch repo.Provider {
	case "gitlab":
	case "github":
	case "":
		return ErrRepositoryMissingProvider
	default:
		return ErrRepositoryInvalidProvider
	}
	return s.store.AddRepository(repo)
}

type Repository struct {
	Name     string
	Owner    string
	Provider string
}

type Pullrequest struct {
	Number int
	URL    string
	Author string
}

func (r *Repository) Equal(other *Repository) bool {
	if r.Name == other.Name && r.Owner == other.Owner && r.Provider == other.Provider {
		return true
	}
	return false
}

func (r *Repository) ToString() string {
	return r.Name + r.Owner + r.Provider
}
