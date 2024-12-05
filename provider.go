package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

type Provider interface {
	GetPullRequests(ctx context.Context, repo Repository, token, base string) ([]*Pullrequest, error)
}

type FakeProvider struct {
	prs map[string][]*Pullrequest
}

func NewFakeProvider() FakeProvider {
	return FakeProvider{prs: make(map[string][]*Pullrequest)}
}

func (f *FakeProvider) GetPullRequests(ctx context.Context, repo Repository, token, base string) ([]*Pullrequest, error) {
	value, success := f.prs[repo.ToString()+token]
	if !success {
		return nil, ErrCannotQueryRepository
	}
	return value, nil
}

func (f *FakeProvider) PullRequestsToReturn(repo Repository, token string, prs []*Pullrequest) {
	f.prs[repo.ToString()+token] = prs
}

type GithubProvider struct {
	username string
}

func NewGithubProvider(username string) *GithubProvider {
	return &GithubProvider{username: username}
}

func (f *GithubProvider) GetPullRequests(ctx context.Context, repo Repository, token, base string) ([]*Pullrequest, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	options := &github.PullRequestListOptions{Base: base, State: "open"}
	prs, _, err := client.PullRequests.List(ctx, repo.Owner, repo.Name, options)
	if err != nil {
		return nil, ErrCannotQueryRepository
	}
	return f.mapToPR(ctx, client, repo.Owner, repo.Name, prs), err
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
		if *details.State == "close" {
			fmt.Printf("Skipping #%d PR due to closed state\n", *pr.Number)
			continue
		}
		if !*details.Mergeable {
			fmt.Printf("Skipping #%d PR due to conflicting state\n", *pr.Number)
			continue
		}
		reviews, _, err := client.PullRequests.ListReviews(ctx, owner, name, *pr.Number, &github.ListOptions{})
		if err != nil {
			fmt.Printf("Skipping #%d PR due to err when listing reviews %v\n", *pr.Number, err)
			continue
		}
		var isApprovedByUser bool
		for _, review := range reviews {
			fmt.Printf("Review %s in state %s", *review.User.Login, *review.State)
			if *review.User.Login == f.username && *review.State == "APPROVED" {
				isApprovedByUser = true
				break
			}
		}
		if isApprovedByUser {
			fmt.Printf("Skipping #%d PR due to already being approved\n", *pr.Number)
			continue
		}
		target = append(target, &Pullrequest{
			Number: *pr.Number,
			URL:    *pr.HTMLURL,
			Author: *pr.User.Login,
			Title:  *pr.Title,
			Opened: pr.CreatedAt.Time,
		})
	}
	return target
}
