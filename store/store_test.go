package store

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	c "tacsiazuma/pullreminder/contract"
)

var sut *SqliteStore

func TestMain(m *testing.M) {
	os.Remove("temp.db")
	db, err := sql.Open("sqlite3", "temp.db")
	if err != nil {
		log.Fatal(err)
	}
	sut = NewSqliteStore(db)
	defer db.Close()
	m.Run()
}

func TestSqliteStore(t *testing.T) {
	t.Parallel()
	t.Run("repositories", func(t *testing.T) {
		t.Parallel()
		t.Run("returns empty slice if no repository added", func(t *testing.T) {
			expected := make([]*c.Repository, 0)
			repos, err := sut.Repositories()
			assert.Nil(t, err, "Should not return error")
			assert.Equal(t, expected, repos, "Should return empty slice")
		})
		t.Run("can add repository", func(t *testing.T) {
			err := sut.AddRepository(&c.Repository{Name: "test", Owner: "me", Provider: "github"})
			assert.NoError(t, err, "Should not return error")
		})
		t.Run("fails to add same repository twice", func(t *testing.T) {
			err := sut.AddRepository(&c.Repository{Name: "test", Owner: "me", Provider: "github"})
			err = sut.AddRepository(&c.Repository{Name: "test", Owner: "me", Provider: "github"})
			assert.Equal(t, c.ErrRepositoryDuplicate, err, "Should return error")
		})
		t.Run("returns previously added repositories", func(t *testing.T) {
			repo := &c.Repository{Name: "test", Owner: "me", Provider: "github"}
			_ = sut.AddRepository(repo)
			repos, err := sut.Repositories()
			assert.Nil(t, err, "Should not return error")
			assert.Equal(t, repo, repos[0], "Should return added repos")
		})
	})
	t.Run("settings", func(t *testing.T) {
		t.Parallel()
		t.Run("can save settings", func(t *testing.T) {
			err := sut.SaveSettings(&c.Settings{ExcludeDraft: true, ExcludeConflicting: true})
			assert.Nil(t, err, "Should not return error")
		})
		t.Run("can get saved settings", func(t *testing.T) {
			expected := &c.Settings{ExcludeDraft: true, ExcludeConflicting: true}
			_ = sut.SaveSettings(expected)
			settings, err := sut.GetSettings()
			assert.Nil(t, err, "Should not return error")
			assert.Equal(t, expected, settings, "Should return same settings")
		})
		t.Run("can update settings", func(t *testing.T) {
			expected := &c.Settings{ExcludeDraft: true, ExcludeConflicting: true}
			_ = sut.SaveSettings(expected)
			expected = &c.Settings{ExcludeDraft: false, ExcludeConflicting: true}
			err := sut.SaveSettings(expected)
			assert.Nil(t, err, "Should not return error")
		})
		t.Run("can get updated settings", func(t *testing.T) {
			expected := &c.Settings{ExcludeDraft: true, ExcludeConflicting: true, Username: "one"}
			_ = sut.SaveSettings(expected)
			expected = &c.Settings{ExcludeDraft: false, ExcludeConflicting: true, Username: "another"}
			_ = sut.SaveSettings(expected)
			settings, err := sut.GetSettings()
			if assert.Nil(t, err, "Should not return error") {
				assert.Equal(t, expected, settings, "Should return updated settings")
			}
		})
	})
	t.Run("credentials", func(t *testing.T) {
		t.Parallel()
		t.Run("returns empty map if no credentials added", func(t *testing.T) {
			expected := make(map[string]string)
			creds, err := sut.Credentials()
			assert.Nil(t, err, "Should not return error")
			assert.Equal(t, expected, creds, "Should return empty map")
		})
		t.Run("can add credentials", func(t *testing.T) {
			err := sut.AddCredentials("github", "token")
			assert.NoError(t, err, "Should not return error")
		})
		t.Run("return previously added credentials", func(t *testing.T) {
			_ = sut.AddCredentials("github", "token")
			creds, err := sut.Credentials()
			assert.Nil(t, err, "Should not return error")
			assert.Equal(t, "token", creds["github"], "Should return token for provider")
		})
	})
}
