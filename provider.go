package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

type Provider interface {
	// Returns the pull requests for a given repository against the provided base branch
	GetPullRequests(ctx context.Context, owner, name, base string) ([]*Pullrequest, error)
}

type FakeProvider struct {
	prs map[string][]*Pullrequest
}

func NewFakeProvider() FakeProvider {
	return FakeProvider{prs: make(map[string][]*Pullrequest)}
}

func (f *FakeProvider) GetPullRequests(ctx context.Context, owner, name, base string) ([]*Pullrequest, error) {
	fmt.Printf("Getting repos for %s/%s\n", owner, name)
	value, success := f.prs[owner+name]
	if !success {
		return nil, ErrCannotQueryRepository
	}
	return value, nil
}

func (f *FakeProvider) PullRequestsToReturn(repo Repository, token string, prs []*Pullrequest) {
	f.prs[repo.Owner+repo.Name] = prs
}

type GithubProvider struct {
	token string
}

func NewGithubProvider(token string) *GithubProvider {
	return &GithubProvider{token: token}
}

func (f *GithubProvider) GetPullRequests(ctx context.Context, owner, repo, base string) ([]*Pullrequest, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: f.token})
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	options := &github.PullRequestListOptions{Base: base}
	prs, _, err := client.PullRequests.List(ctx, owner, repo, options)
	if err != nil {
		return nil, ErrCannotQueryRepository
	}
	return f.mapToPR(ctx, client, owner, repo, prs), err
}

func (f *GithubProvider) mapToPR(ctx context.Context, client *github.Client, owner, name string, origin []*github.PullRequest) (target []*Pullrequest) {
	target = make([]*Pullrequest, 0)
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
		target = append(target, &Pullrequest{
			Number:      *pr.Number,
			URL:         *pr.HTMLURL,
			Author:      *pr.User.Login,
			Title:       *pr.Title,
			Opened:      pr.CreatedAt.Time,
			Assignee:    assignee,
			Description: description,
			Mergeable:   *details.Mergeable,
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

func MapReviews(reviews []*github.PullRequestReview) []Review {
	result := make([]Review, len(reviews))
	for i, r := range reviews {
		result[i] = Review{Author: *r.User.Login, Body: *r.Body, State: *r.State}
	}
	return result
}
