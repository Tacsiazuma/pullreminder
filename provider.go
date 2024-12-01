package main

import (
	"context"

	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

type Provider interface {
	GetPullRequests(ctx context.Context, repo Repository, token string) ([]*Pullrequest, error)
}

type FakeProvider struct {
	prs map[string][]*Pullrequest
}

func NewFakeProvider() FakeProvider {
	return FakeProvider{prs: make(map[string][]*Pullrequest)}
}

func (f *FakeProvider) GetPullRequests(ctx context.Context, repo Repository, token string) ([]*Pullrequest, error) {
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
}

func NewGithubProvider() *GithubProvider {
	return &GithubProvider{}
}

func (f *GithubProvider) GetPullRequests(ctx context.Context, repo Repository, token string) ([]*Pullrequest, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	prs, _, err := client.PullRequests.List(ctx, repo.Owner, repo.Name, nil)
	if err != nil {
		return nil, ErrCannotQueryRepository
	}
	return mapToPR(prs), err
}

func mapToPR(origin []*github.PullRequest) (target []*Pullrequest) {
	target = make([]*Pullrequest, len(origin))
	for i := range origin {
		target[i] = &Pullrequest{}
	}
	return target
}
