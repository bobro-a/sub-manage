package app

import (
	"context"
	"fmt"
	"net/http"
	"sub-manage/internal/handler"
	"sub-manage/internal/repo"
	"sub-manage/internal/usecase"
	"sub-manage/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type App struct {
	ctx context.Context
	cfg *config.Config
	db  *sqlx.DB
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	db, err := sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database is unreachable: %w", err)
	}
	return &App{
		ctx: ctx,
		cfg: cfg,
		db:  db,
	}, nil
}

func (a *App) Start() error {
	defer func() {
		_ = a.db.Close()
	}()

	m, err := migrate.New(a.cfg.Migrations.Path, a.cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to init migration: %w", err)
	}
	if err = m.Up(); err != nil && err.Error() != "no change" {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	subRepo := repo.NewSubRepo(a.cfg.Database.Name, a.db)
	subUseCase := usecase.New(subRepo)
	subHandler := handler.NewSubHandler(subUseCase)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			subHandler.Create(w, r)
		case http.MethodGet:
			if r.URL.Query().Has("id") {
				subHandler.Read(w, r)
			} else {
				subHandler.List(w, r)
			}
		case http.MethodDelete:
			subHandler.Delete(w, r)
		case http.MethodPut:
			subHandler.Update(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/sum", subHandler.Sum)
	_ = http.ListenAndServe(":9090", nil)
	return nil
}
