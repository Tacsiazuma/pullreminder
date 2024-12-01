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
	assert.NotEqual(t, token, "", "A valid github token should be set!")
	sut := NewGithubProvider()
	t.Run("instance", func(t *testing.T) {
		t.Run("can be instantiated without a token", func(t *testing.T) {
			assert.NotNil(t, sut, "should not be nil")
		})
	})
	t.Run("PullRequests Function", func(t *testing.T) {
		t.Run("fail if the provided token is not valid", func(t *testing.T) {
			repo := &Repository{Name: "pullreminder-test", Owner: "tacsiazuma"}
			prs, err := sut.GetPullRequests(ctx, *repo, "invalid token")
			assert.Nil(t, prs, "Should not return pull requests")
			assert.Equal(t, ErrCannotQueryRepository, err, "Should return error")
		})
		t.Run("return empty slice if no PRs", func(t *testing.T) {
			repo := &Repository{Name: "pullreminder-test", Owner: "tacsiazuma"}
			expected := make([]*Pullrequest, 0)
			prs, err := sut.GetPullRequests(ctx, *repo, token)
			assert.Equal(t, prs, expected, "Should not return pull requests")
			assert.Nil(t, err, "Should return not error")
		})
	})
}
