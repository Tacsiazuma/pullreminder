package store

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	c "tacsiazuma/pullreminder/contract"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(db *sql.DB) *SqliteStore {
	_, err := db.Exec("create table if not exists repositories (name varchar, owner varchar, provider varchar, primary key (name, owner,provider))")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("create table if not exists credentials (provider varchar, token varchar)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("create table if not exists settings (key varchar, value text, primary key (key))")
	if err != nil {
		log.Fatal(err)
	}
	return &SqliteStore{db: db}
}

func (s *SqliteStore) AddRepository(repo *c.Repository) error {
	_, err := s.db.Exec("insert into repositories (name, owner, provider) values (?,?,?)", repo.Name, repo.Owner, repo.Provider)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint") {
		return c.ErrRepositoryDuplicate
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

func (s *SqliteStore) SaveSettings(settings *c.Settings) error {
	marshalled, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?,?)", "settings", marshalled)
	return err
}

func (s *SqliteStore) GetSettings() (*c.Settings, error) {
	rows, err := s.db.Query("SELECT value FROM settings where key = 'settings'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	settings := &c.Settings{}
	for rows.Next() {
		var marshalled string
		err = rows.Scan(&marshalled)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(marshalled), settings)
		if err != nil {
			return nil, err
		}
	}
	return settings, nil
}

func (s *SqliteStore) Repositories() ([]*c.Repository, error) {
	rows, err := s.db.Query("SELECT name, owner, provider FROM repositories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	repos := make([]*c.Repository, 0)
	for rows.Next() {
		var name string
		var owner string
		var provider string
		err = rows.Scan(&name, &owner, &provider)
		if err != nil {
			return nil, err
		}
		repos = append(repos, &c.Repository{Name: name, Owner: owner, Provider: provider})
	}
	return repos, nil
}
