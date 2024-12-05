package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.TODO()
	token := os.Getenv("GITHUB_TOKEN")
	username := "Tacsiazuma"
	assert.NotEqual(t, token, "", "A valid github token should be set!")
	sut := NewGithubProvider(username)
	t.Run("instance", func(t *testing.T) {
		t.Run("can be instantiated without a token", func(t *testing.T) {
			assert.NotNil(t, sut, "should not be nil")
		})
	})
	t.Run("PullRequests Function", func(t *testing.T) {
		t.Run("fail if the provided token is not valid", func(t *testing.T) {
			repo := &Repository{Name: "pullreminder-test", Owner: "tacsiazuma"}
			prs, err := sut.GetPullRequests(ctx, *repo, "invalid token", "master")
			assert.Nil(t, prs, "Should not return pull requests")
			assert.Equal(t, ErrCannotQueryRepository, err, "Should return error")
		})
		t.Run("return empty slice if no PRs on the base branch", func(t *testing.T) {
			repo := &Repository{Name: "pullreminder-test", Owner: "tacsiazuma"}
			expected := make([]*Pullrequest, 0)
			prs, err := sut.GetPullRequests(ctx, *repo, token, "master")
			assert.Equal(t, prs, expected, "Should not return pull requests")
			assert.Nil(t, err, "Should not return error")
		})
		t.Run("return open PRs opened against the base branch", func(t *testing.T) {
			repo := &Repository{Name: "pullreminder-test", Owner: "tacsiazuma"}
			prs, err := sut.GetPullRequests(ctx, *repo, token, "main")
			expected := &Pullrequest{Number: 1, URL: "https://github.com/Tacsiazuma/pullreminder-test/pull/1"}
			if assert.NotNil(t, prs) && assert.Equal(t, 1, len(prs), "Should contain pull requests") {
				assert.Equal(t, prs[0], expected, "Should return pull requests")
				assert.Nil(t, err, "Should not return error")
			}
		})
		t.Run("does not return PRs approved by user", func(t *testing.T) {
			repo := &Repository{Name: "pullreminder-test", Owner: "tacsiazuma"}
			prs, err := sut.GetPullRequests(ctx, *repo, token, "approved")
			assert.Equal(t, 0, len(prs), "Should not return PRs")
			assert.Nil(t, err, "Should not return error")
		})
		// https://github.com/Tacsiazuma/pullreminder-test/pull/3 approved-pr -> closed
		t.Run("does not return closed PRs", func(t *testing.T) {
			repo := &Repository{Name: "pullreminder-test", Owner: "tacsiazuma"}
			prs, err := sut.GetPullRequests(ctx, *repo, token, "closed")
			assert.Equal(t, 0, len(prs), "Should not return closed PRs")
			assert.Nil(t, err, "Should not return error")
		})
	})
}
