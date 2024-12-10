package provider

import (
	"context"
	"errors"
	"fmt"
	"tacsiazuma/pullreminder/contract"

	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

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

type GithubProvider struct {
	token string
}

func NewGithubProvider(token string) *GithubProvider {
	return &GithubProvider{token: token}
}

func (f *GithubProvider) GetPullRequests(ctx context.Context, owner, repo, base string) ([]*contract.Pullrequest, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: f.token})
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	options := &github.PullRequestListOptions{Base: base}
	prs, _, err := client.PullRequests.List(ctx, owner, repo, options)
	if err != nil {
		return nil, ErrCannotQueryRepository
	}
	return f.mapToPR(ctx, client, owner, repo, prs), err
}

func (f *GithubProvider) mapToPR(ctx context.Context, client *github.Client, owner, name string, origin []*github.PullRequest) (target []*contract.Pullrequest) {
	target = make([]*contract.Pullrequest, 0)
	for _, pr := range origin {
		fmt.Printf("Processing #%d PR\n", *pr.Number)
		details, _, err := client.PullRequests.Get(ctx, owner, name, *pr.Number)
		if err != nil {
			fmt.Printf("Skipping #%d PR due to err %v\n", *pr.Number, err)
			continue
		}
		reviews, _, err := client.PullRequests.ListReviews(ctx, owner, name, *pr.Number, &github.ListOptions{})
		if err != nil {
			fmt.Printf("Skipping #%d PR due to err when listing reviews %v\n", *pr.Number, err)
			continue
		}
		var description string
		if pr.Body == nil {
			description = ""
		} else {
			description = *pr.Body
		}
		var assignee string
		if pr.Assignee == nil {
			assignee = ""
		} else {
			assignee = *pr.Assignee.Login
		}
		var mergeable bool
		if details.Mergeable == nil || !*details.Mergeable {
			mergeable = false
		} else {
			mergeable = true
		}
		target = append(target, &contract.Pullrequest{
			Number:      *pr.Number,
			URL:         *pr.HTMLURL,
			Author:      *pr.User.Login,
			Title:       *pr.Title,
			Opened:      pr.CreatedAt.Time,
			Assignee:    assignee,
			Description: description,
			Mergeable:   mergeable,
			Reviewers:   MapReviewers(pr.RequestedReviewers),
			Reviews:     MapReviews(reviews),
		})
	}
	return target
}

func MapReviewers(reviewers []*github.User) []string {
	result := make([]string, len(reviewers))
	for i, r := range reviewers {
		result[i] = *r.Login
	}
	return result
}

func MapReviews(reviews []*github.PullRequestReview) []contract.Review {
	result := make([]contract.Review, len(reviews))
	for i, r := range reviews {
		result[i] = contract.Review{Author: *r.User.Login, Body: *r.Body, State: *r.State}
	}
	return result
}
