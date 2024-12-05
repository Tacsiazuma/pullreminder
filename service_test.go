package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	t.Run("can be instantiated", func(t *testing.T) {
		_ = CreateService()
	})
	t.Run("NeedsAttention function", func(t *testing.T) {
		ctx := context.TODO()
		t.Run("error", func(t *testing.T) {
			t.Run("when no repositories provided", func(t *testing.T) {
				sut := CreateService()
				result, err := sut.NeedsAttention(ctx)
				assert.Equal(t, ErrNoRepositoriesProvided, err, "Should return error")
				assert.Nil(t, result, "Should return nil")
			})
			t.Run("when no credentials provided", func(t *testing.T) {
				sut := CreateService()
				repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
				_ = sut.AddRepository(repo)
				_, err := sut.NeedsAttention(ctx)
				assert.Equal(t, ErrNoCredentialsProvidedForGithub, err, "Should return error")
			})
			t.Run("when different providers credentials added", func(t *testing.T) {
				sut := CreateService()
				_ = sut.AddCredentials("github", "sometoken")
				repo := &Repository{Name: "reponame", Owner: "owner", Provider: "gitlab"}
				_ = sut.AddRepository(repo)
				repos, err := sut.NeedsAttention(ctx)
				assert.Nil(t, repos, "Should not return result")
				assert.Equal(t, ErrNoCredentialsProvidedForGitlab, err, "Should return error")
			})
			t.Run("when provider is failing to query repo", func(t *testing.T) {
				sut := CreateService()
				repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
				_ = sut.AddCredentials("github", "sometoken")
				_ = sut.AddRepository(repo)
				prs, err := sut.NeedsAttention(ctx)
				assert.Equal(t, make([]*Pullrequest, 0), prs, "Should return empty list")
				assert.Equal(t, ErrCannotQueryRepository, err, "Should return error")
			})
		})
		t.Run("empty list when provider returns empty list", func(t *testing.T) {
			sut := CreateService()
			repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
			expected := make([]*Pullrequest, 0)
			provider.PullRequestsToReturn(*repo, "sometoken", expected)
			_ = sut.AddCredentials("github", "sometoken")
			_ = sut.AddRepository(repo)
			prs, err := sut.NeedsAttention(ctx)
			assert.Equal(t, expected, prs, "Should return empty list")
			assert.Nil(t, err, "Should not return error")
		})
		t.Run("does not return conflicting PRs", func(t *testing.T) {
			sut := CreateService()
			repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
			_ = sut.AddCredentials("github", "sometoken")
			_ = sut.AddRepository(repo)
			expected := CreateConflictingPR()
			provider.PullRequestsToReturn(*repo, "sometoken", expected)
			prs, err := sut.NeedsAttention(ctx)
			assert.Equal(t, 0, len(prs), "Should return empty list")
			assert.Nil(t, err, "Should not return error")
		})
		t.Run("does not return accepted PRs by the user", func(t *testing.T) {
			sut := CreateService()
			repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
			_ = sut.AddCredentials("github", "sometoken")
			_ = sut.AddRepository(repo)
			expected := CreateConflictingPR()
			provider.PullRequestsToReturn(*repo, "sometoken", expected)
			prs, err := sut.NeedsAttention(ctx)
			assert.Equal(t, 0, len(prs), "Should return empty list")
			assert.Nil(t, err, "Should not return error")
		})
	})
	t.Run("AddRepository function", func(t *testing.T) {
		sut := CreateService()
		t.Run("returns error when no name provided", func(t *testing.T) {
			repo := &Repository{}
			err := sut.AddRepository(repo)
			assert.Equal(t, ErrRepositoryMissingName, err, "Should return error")
		})
		t.Run("returns error when no owner provided", func(t *testing.T) {
			repo := &Repository{Name: "reponame"}
			err := sut.AddRepository(repo)
			assert.Equal(t, ErrRepositoryMissingOwner, err, "Should return error")
		})
		t.Run("returns error when no provider passed", func(t *testing.T) {
			repo := &Repository{Name: "reponame", Owner: "owner"}
			err := sut.AddRepository(repo)
			assert.Equal(t, ErrRepositoryMissingProvider, err, "Should return error")
		})
		t.Run("returns error when invalid provider passed", func(t *testing.T) {
			repo := &Repository{Name: "reponame", Owner: "owner", Provider: "notgithub"}
			err := sut.AddRepository(repo)
			assert.Equal(t, ErrRepositoryInvalidProvider, err, "Should return error")
		})
		t.Run("does not return error when valid provider passed", func(t *testing.T) {
			repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
			err := sut.AddRepository(repo)
			assert.Nil(t, err, "Should not return error")
		})
		t.Run("returns error when adding same repository twice", func(t *testing.T) {
			sut := CreateService()
			repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
			_ = sut.AddRepository(repo)
			err := sut.AddRepository(repo)
			assert.Equal(t, ErrRepositoryDuplicate, err, "Should return error")
		})
	})
	t.Run("Repositories function", func(t *testing.T) {
		t.Run("returns and empty slice if no repository added", func(t *testing.T) {
			sut := CreateService()
			list, err := sut.Repositories()
			assert.Empty(t, list, "It should return an empty slice")
			assert.Nil(t, err, "Should not return error")
		})
		t.Run("returns repository after adding it", func(t *testing.T) {
			sut := CreateService()
			repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
			_ = sut.AddRepository(repo)
			repos, err := sut.Repositories()
			assert.Nil(t, err, "Should not return error")
			assert.Equal(t, repo, repos[0], "Should return the same")
		})
	})
	t.Run("AddCredentials function", func(t *testing.T) {
		t.Run("fails when invoked with invalid provider", func(t *testing.T) {
			sut := CreateService()
			err := sut.AddCredentials("notgithub", "sometoken")
			assert.Equal(t, ErrInvalidProvider, err, "Should return error")
		})
		t.Run("succeeds when invoked with valid provider", func(t *testing.T) {
			sut := CreateService()
			err := sut.AddCredentials("github", "sometoken")
			assert.Nil(t, err, "Should not return error")
		})
	})
}

var provider FakeProvider

func CreateService() *Service {
	provider = NewFakeProvider()
	return New(&provider)
}

func CreateConflictingPR() []*Pullrequest {
	return []*Pullrequest{&Pullrequest{}}
}

// feed in repositories and credentials
// get list of PRs which requires attention
// invoke provider to check if repositories are available
// unavailable repositories will be marked with error
// invoke provider to get list of open PRs associated to repositories
// return PRs or empty slice
// ability to check draft PRs as well
// ask if there is anything to do
// list repositories fed in
