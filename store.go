package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

type Store interface {
	AddRepository(repo *Repository) error
	AddCredentials(provider, token string) error
	Repositories() ([]*Repository, error)
	Credentials() (map[string]string, error)
}

type FakeStore struct {
	repositories []*Repository
	credentails  map[string]string
}

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(db *sql.DB) Store {
	_, err := db.Exec("create table if not exists repositories (name varchar, owner varchar, provider varchar, primary key (name, owner,provider))")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("create table if not exists credentials (provider varchar, token varchar)")
	if err != nil {
		log.Fatal(err)
	}
	return &SqliteStore{db: db}
}

func (s *SqliteStore) AddRepository(repo *Repository) error {
	_, err := s.db.Exec("insert into repositories (name, owner, provider) values (?,?,?)", repo.Name, repo.Owner, repo.Provider)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint") {
		return ErrRepositoryDuplicate
	}
	return nil
}

func (s *SqliteStore) AddCredentials(provider, token string) error {
	_, err := s.db.Exec("insert into credentials (provider, token) values (?,?)", provider, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteStore) Credentials() (map[string]string, error) {
	rows, err := s.db.Query("SELECT provider, token FROM credentials")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	creds := make(map[string]string)
	for rows.Next() {
		var token string
		var provider string
		err = rows.Scan(&provider, &token)
		if err != nil {
			return nil, err
		}
		creds[provider] = token
	}
	return creds, nil
}

func (s *SqliteStore) Repositories() ([]*Repository, error) {
	rows, err := s.db.Query("SELECT name, owner, provider FROM repositories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	repos := make([]*Repository, 0)
	for rows.Next() {
		var name string
		var owner string
		var provider string
		err = rows.Scan(&name, &owner, &provider)
		if err != nil {
			return nil, err
		}
		repos = append(repos, &Repository{Name: name, Owner: owner, Provider: provider})
	}
	return repos, nil
}
func NewFakeStore() Store {
	return &FakeStore{repositories: make([]*Repository, 0), credentails: make(map[string]string)}
}

func (s *FakeStore) AddRepository(repo *Repository) error {
	for _, r := range s.repositories {
		if r.Equal(repo) {
			return ErrRepositoryDuplicate
		}
	}
	s.repositories = append(s.repositories, repo)
	return nil
}

func (s *FakeStore) AddCredentials(provider, token string) error {
	s.credentails[provider] = token
	return nil
}

func (s *FakeStore) Credentials() (map[string]string, error) {
	return s.credentails, nil
}

func (s *FakeStore) Repositories() ([]*Repository, error) {
	return s.repositories, nil
}
