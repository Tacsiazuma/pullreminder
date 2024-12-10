package contract

import (
	"errors"
	"time"
)

type Pullrequest struct {
	Number      int
	URL         string
	Author      string
	Title       string
	Opened      time.Time
	Assignee    string
	Reviewers   []string
	Description string
	Mergeable   bool
	Draft       bool
	Reviews     []Review
}

type Review struct {
	Body   string
	State  string
	Author string
}

type Repository struct {
	Name     string
	Owner    string
	Provider string
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

func (r *Repository) Equal(other *Repository) bool {
	if r.Name == other.Name && r.Owner == other.Owner && r.Provider == other.Provider {
		return true
	}
	return false
}

func (r *Repository) ToString() string {
	return r.Name + r.Owner + r.Provider
}

type Settings struct {
	ExcludeDraft       bool
	ExcludeConflicting bool
}
