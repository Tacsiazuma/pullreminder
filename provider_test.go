package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestGithubProvider(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.TODO()
	token := os.Getenv("GITHUB_TOKEN")
	username := "Tacsiazuma"
	assert.NotEqual(t, token, "", "A valid github token should be set!")
	sut := NewGithubProvider(username, token)
	t.Run("instance", func(t *testing.T) {
		t.Run("can be instantiated without a token", func(t *testing.T) {
			assert.NotNil(t, sut, "should not be nil")
		})
	})
	t.Run("PullRequests Function", func(t *testing.T) {
		reponame := "pullreminder-test"
		owner := "tacsiazuma"
		t.Run("fail if the provided token is not valid", func(t *testing.T) {
			sut := NewGithubProvider(username, "invalid token")
			prs, err := sut.GetPullRequests(ctx, reponame, owner, "master")
			assert.Nil(t, prs, "Should not return pull requests")
			assert.Equal(t, ErrCannotQueryRepository, err, "Should return error")
		})
		t.Run("return empty slice if no PRs on the base branch", func(t *testing.T) {
			expected := make([]*Pullrequest, 0)
			prs, err := sut.GetPullRequests(ctx, reponame, owner, "master")
			assert.Equal(t, prs, expected, "Should not return pull requests")
			assert.Nil(t, err, "Should not return error")
		})
		t.Run("return open PRs opened against the base branch", func(t *testing.T) {
			prs, err := sut.GetPullRequests(ctx, reponame, owner, "main")
			rev := &Review{Author: "Tacsiazuma", Body: "Oh noes", State: "COMMENTED"}
			expected := &Pullrequest{Number: 1,
				URL:         "https://github.com/Tacsiazuma/pullreminder-test/pull/1",
				Author:      "Tacsiazuma",
				Title:       "Update LICENSE",
				Description: "description provided",
				Opened:      time.Date(2024, 12, 1, 22, 3, 25, 0, time.UTC),
				Assignee:    "Tacsiazuma",
				Mergeable:   true,
				Reviewers:   []string{"letscodehu"},
				Reviews:     []Review{*rev},
			}
			if assert.NotNil(t, prs) && assert.Equal(t, 1, len(prs), "Should contain pull requests") {
				assert.Equal(t, expected, prs[0], "Should return open PRs")
				assert.Nil(t, err, "Should not return error")
			}
		})
		t.Run("does return PRs approved by user", func(t *testing.T) {
			prs, err := sut.GetPullRequests(ctx, reponame, owner, "approved")
			assert.Equal(t, 1, len(prs), "Should return approved PRs")
			assert.Nil(t, err, "Should not return error")
		})
		// https://github.com/Tacsiazuma/pullreminder-test/pull/3 approved-pr -> closed
		t.Run("does not return closed PRs", func(t *testing.T) {
			prs, err := sut.GetPullRequests(ctx, reponame, owner, "closed")
			assert.Equal(t, 0, len(prs), "Should not return closed PRs")
			assert.Nil(t, err, "Should not return error")
		})
		t.Run("does return conflicting PRs", func(t *testing.T) {
			prs, err := sut.GetPullRequests(ctx, reponame, owner, "conflicting")
			assert.Equal(t, 1, len(prs), "Should return conflicting PRs")
			assert.Nil(t, err, "Should not return error")
		})
	})
}
