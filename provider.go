package main

type Provider interface {
	GetPullRequests(repo Repository, token string) ([]*Pullrequest, error)
}

type FakeProvider struct {
	prs map[string][]*Pullrequest
}

func NewFakeProvider() FakeProvider {
	return FakeProvider{prs: make(map[string][]*Pullrequest)}
}

func (f *FakeProvider) GetPullRequests(repo Repository, token string) ([]*Pullrequest, error) {
	value, success := f.prs[repo.ToString()+token]
	if !success {
		return nil, ErrCannotQueryRepository
	}
	return value, nil
}

func (f *FakeProvider) PullRequestsToReturn(repo Repository, token string, prs []*Pullrequest) {
	f.prs[repo.ToString()+token] = prs
}
