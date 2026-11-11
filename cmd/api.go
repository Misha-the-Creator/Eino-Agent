package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Misha-the-Creator/Eino-Agent/generated"
	"github.com/Misha-the-Creator/Eino-Agent/internal/agent"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Logger    *slog.Logger
	pool      *pgxpool.Pool
	addr      string
	apiKey    string
	modelName string
	ctx       context.Context
}

func (app *App) mount() (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//
	
	q := generated.New(app.pool)
	MAS, err := BuildMAS(app.ctx, app, q)

	if err != nil {
		app.Logger.Error("Ошибка при сборке MAS", "err", err)
		return nil, err
	}

	agentSvc := agent.NewSvc(app.Logger)
	agentHndlr := agent.NewHndlr(agentSvc, app.Logger, MAS)
	r.Post("/agent", agentHndlr.AgentRoute)

	return r, nil
}

func (app *App) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.addr,
		Handler:      h,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,
	}

	app.Logger.Info("Сервер запущен на адресе", "addr", app.addr)

	return srv.ListenAndServe()
}
