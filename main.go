package main

import "errors"

func main() {
}

type Service struct {
	repositories []*Repository
	credentails  map[string]string
}

func New() *Service {
	return &Service{repositories: make([]*Repository, 0), credentails: make(map[string]string)}
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
)

func (s *Service) NeedsAttention() (any, error) {
	for _, repo := range s.repositories {
		if s.credentails[repo.Provider] == "" {
			switch repo.Provider {
			case "gitlab":
				return nil, ErrNoCredentialsProvidedForGitlab
			case "github":
				return nil, ErrNoCredentialsProvidedForGithub
			}
		}
	}
	return nil, ErrNoRepositoriesProvided
}

func (s *Service) Repositories() ([]*Repository, error) {
	return s.repositories, nil
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
	s.credentails[provider] = token
	return nil
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
	for _, r := range s.repositories {
		if r.Equal(repo) {
			return ErrRepositoryDuplicate
		}
	}
	s.repositories = append(s.repositories, repo)
	return nil
}

type Repository struct {
	Name     string
	Owner    string
	Provider string
}

func (r *Repository) Equal(other *Repository) bool {
	if r.Name == other.Name && r.Owner == other.Owner && r.Provider == other.Provider {
		return true
	}
	return false
}
