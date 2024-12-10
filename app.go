package main

import (
	"context"
	"database/sql"
	"github.com/gen2brain/beeep"
	"github.com/go-co-op/gocron/v2"
	"github.com/kirsle/configdir"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
	"path/filepath"
	"sync"
	c "tacsiazuma/pullreminder/contract"
	p "tacsiazuma/pullreminder/provider"
	s "tacsiazuma/pullreminder/service"
	st "tacsiazuma/pullreminder/store"
)

// App struct
type App struct {
	ctx       context.Context
	service   *s.Service
	prs       []*c.Pullrequest
	m         sync.RWMutex
	scheduler gocron.Scheduler
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

var db *sql.DB

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	configpath := configdir.LocalConfig("pullreminder")
	err := configdir.MakePath(configpath)
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("sqlite3", filepath.Join(configpath, "db.sqlite"))
	if err != nil {
		log.Fatal(err)
	}
	a.service = s.NewService(p.NewGithubProvider(os.Getenv("GITHUB_TOKEN")), st.NewSqliteStore(db))
	a.ctx = ctx
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}
	a.scheduler = s
	a.scheduler.Start()
}

func (a *App) UpdateSchedule(cron string) {
	a.m.Lock()
	defer a.m.Unlock()
	a.scheduler.Shutdown()
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}
	a.scheduler = s
	runtime.LogPrintf(a.ctx, "Notification schedule set to %s", cron)
	_, err = a.scheduler.NewJob(
		gocron.CronJob(cron, false), gocron.NewTask(
			func() {
				runtime.LogPrint(a.ctx, "notification sent")
				prs, err := a.CheckPRs()
				if err != nil {
					runtime.LogPrint(a.ctx, err.Error())
				}
				if len(prs) > 0 {
					_ = beeep.Notify("Hey you!", "A PR requires your attention!", "frontend/src/assets/images/logo-universal.png")
				}
			},
		))
	if err != nil {
		log.Fatal(err)
	}
	a.scheduler.Start()
}

func (a *App) OnShutdown(ctx context.Context) {
	_ = a.scheduler.Shutdown()
	db.Close()
}

// repos returns the list of repositories
func (a *App) Repos() ([]*c.Repository, error) {
	return a.service.Repositories()
}

// repos returns the list of repositories
func (a *App) AddRepo(repo *c.Repository) error {
	return a.service.AddRepository(repo)
}

// repos returns the list of repositories
func (a *App) CheckPRs() ([]*c.Pullrequest, error) {
	a.m.Lock()
	defer a.m.Unlock()
	if a.prs == nil {
		prs, err := a.service.NeedsAttention(a.ctx)
		if err != nil {
			return nil, err
		}
		a.prs = prs
	}
	return a.prs, nil
}
