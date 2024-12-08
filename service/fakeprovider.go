package service

import (
	"context"
	"fmt"
	c "tacsiazuma/pullreminder/contract"
)

type FakeProvider struct {
	prs map[string][]*c.Pullrequest
}

func NewFakeProvider() FakeProvider {
	return FakeProvider{prs: make(map[string][]*c.Pullrequest)}
}

func (f *FakeProvider) GetPullRequests(ctx context.Context, owner, name, base string) ([]*c.Pullrequest, error) {
	fmt.Printf("Getting repos for %s/%s\n", owner, name)
	value, success := f.prs[owner+name]
	if !success {
		return nil, c.ErrCannotQueryRepository
	}
	return value, nil
}

func (f *FakeProvider) PullRequestsToReturn(repo c.Repository, token string, prs []*c.Pullrequest) {
	f.prs[repo.Owner+repo.Name] = prs
}
