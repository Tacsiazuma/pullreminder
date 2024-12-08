package service

import (
	"context"
	c "tacsiazuma/pullreminder/contract"
)

type Service struct {
	store    Store
	provider Provider
}

type FakeStore struct {
	repositories []*c.Repository
	credentails  map[string]string
}

func New(provider Provider, store Store) *Service {
	return &Service{
		store:    store,
		provider: provider,
	}
}

func (s *Service) NeedsAttention(ctx context.Context) ([]*c.Pullrequest, error) {
	repos, _ := s.store.Repositories()
	if len(repos) == 0 {
		return nil, c.ErrNoRepositoriesProvided
	}
	creds, _ := s.store.Credentials()
	total := make([]*c.Pullrequest, 0)
	for _, repo := range repos {
		if creds[repo.Provider] == "" {
			switch repo.Provider {
			case "gitlab":
				return nil, c.ErrNoCredentialsProvidedForGitlab
			case "github":
				return nil, c.ErrNoCredentialsProvidedForGithub
			}
		}
		prs, err := s.provider.GetPullRequests(ctx, repo.Owner, repo.Name, "main")
		if err != nil {
			return total, err
		}
		total = append(total, s.filter(prs)...)
	}
	return total, nil
}

func (s *Service) filter(origin []*c.Pullrequest) []*c.Pullrequest {
	var result []*c.Pullrequest
	for _, pr := range origin {
		if !pr.Mergeable {
			continue
		}
		result = append(result, pr)
	}
	return result
}

func (s *Service) Repositories() ([]*c.Repository, error) {
	return s.store.Repositories()
}

func (s *Service) AddCredentials(provider, token string) error {
	switch provider {
	case "github":
		break
	case "gitlab":
		break
	default:
		return c.ErrInvalidProvider
	}
	return s.store.AddCredentials(provider, token)
}

func (s *Service) AddRepository(repo *c.Repository) error {
	if repo.Name == "" {
		return c.ErrRepositoryMissingName
	}
	if repo.Owner == "" {
		return c.ErrRepositoryMissingOwner
	}
	switch repo.Provider {
	case "gitlab":
	case "github":
	case "":
		return c.ErrRepositoryMissingProvider
	default:
		return c.ErrRepositoryInvalidProvider
	}
	return s.store.AddRepository(repo)
}


type Store interface {
	AddRepository(repo *c.Repository) error
	AddCredentials(provider, token string) error
	Repositories() ([]*c.Repository, error)
	Credentials() (map[string]string, error)
}

type Provider interface {
	// Returns the pull requests for a given repository against the provided base branch
	GetPullRequests(ctx context.Context, owner, name, base string) ([]*c.Pullrequest, error)
}
