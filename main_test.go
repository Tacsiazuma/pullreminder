package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItWorks(t *testing.T) {
	_ = CreateService()
}

func TestItReturnsErrorWhenNoRepositoriesProvided(t *testing.T) {
	sut := CreateService()
	result, err := sut.NeedsAttention()
	assert.Equal(t, ErrNoRepositoriesProvided, err, "Should return error")
	assert.Nil(t, result, "Should return nil")
}

func TestItReturnsEmptySliceIfNoRepositoriesFedIn(t *testing.T) {
	sut := CreateService()
	list, err := sut.Repositories()
	assert.Empty(t, list, "It should return an empty slice")
	assert.Nil(t, err, "Should not return error")
}

func TestItReturnsErrorWhenRepositoryWithoutNameProvided(t *testing.T) {
	sut := CreateService()
	repo := &Repository{}
	err := sut.AddRepository(repo)
	assert.Equal(t, ErrRepositoryMissingName, err, "Should return error")
}

func TestItReturnsErrorWhenRepositoryWithoutOwnerProvided(t *testing.T) {
	sut := CreateService()
	repo := &Repository{Name: "reponame"}
	err := sut.AddRepository(repo)
	assert.Equal(t, ErrRepositoryMissingOwner, err, "Should return error")
}

func TestItReturnsErrorWhenRepositoryWithoutProviderPassed(t *testing.T) {
	sut := CreateService()
	repo := &Repository{Name: "reponame", Owner: "owner"}
	err := sut.AddRepository(repo)
	assert.Equal(t, ErrRepositoryMissingProvider, err, "Should return error")
}

func TestItReturnsErrorWhenRepositoryWithInvalidProvider(t *testing.T) {
	sut := CreateService()
	repo := &Repository{Name: "reponame", Owner: "owner", Provider: "notgithub"}
	err := sut.AddRepository(repo)
	assert.Equal(t, ErrRepositoryInvalidProvider, err, "Should return error")
}

func TestItDoesNotReturnErrorWhenValidRepositoryProvided(t *testing.T) {
	sut := CreateService()
	repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
	err := sut.AddRepository(repo)
	assert.Nil(t, err, "Should not return error")
}

func TestItReturnsRepositoryAfterAddingIt(t *testing.T) {
	sut := CreateService()
	repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
	_ = sut.AddRepository(repo)
	repos, err := sut.Repositories()
	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, repo, repos[0], "Should return the same")
}

func TestItReturnsErrorWhenAddingSameRepoTwice(t *testing.T) {
	sut := CreateService()
	repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
	_ = sut.AddRepository(repo)
	err := sut.AddRepository(repo)
	assert.Equal(t, ErrRepositoryDuplicate, err, "Should return error")
}

func TestNeedsAttentionReturnsErrorWhenNoCredentialsProvided(t *testing.T) {
	sut := CreateService()
	repo := &Repository{Name: "reponame", Owner: "owner", Provider: "github"}
	_ = sut.AddRepository(repo)
	_, err := sut.NeedsAttention()
	assert.Equal(t, ErrNoCredentialsProvidedForGithub, err, "Should return error")
}

func TestAddCredentialsFailsWhenUsedWithInvalidProvider(t *testing.T) {
	sut := CreateService()
	err := sut.AddCredentials("notgithub", "sometoken")
	assert.Equal(t, ErrInvalidProvider, err, "Should return error")
}

func TestAddCredentialsSucceedsWhenUsedWithValidProvider(t *testing.T) {
	sut := CreateService()
	err := sut.AddCredentials("github", "sometoken")
	assert.Nil(t, err, "Should not return error")
}

func TestNeedsAttentionReturnsErrorWhenDifferentProviderCredentialProvided(t *testing.T) {
	sut := CreateService()
	_ = sut.AddCredentials("github", "sometoken")
	repo := &Repository{Name: "reponame", Owner: "owner", Provider: "gitlab"}
	_ = sut.AddRepository(repo)
	repos, err := sut.NeedsAttention()
	assert.Nil(t, repos, "Should not return result")
	assert.Equal(t, ErrNoCredentialsProvidedForGitlab, err, "Should return error")
}

func CreateService() *Service {
	return New()
}

// feed in repositories and credentials
// ask if there is anything to do
// list repositories fed in
