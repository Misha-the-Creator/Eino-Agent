package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Misha-the-Creator/Eino-Agent/internal/agent"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Logger      *slog.Logger
	pool        *pgxpool.Pool
	addr        string
}

func (app *App) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// forecastStore := forecasttable.NewStore(app.Logger)
	// loadFactRepo := factloading.NewRepo(queries, app.Logger)
	// loadFactSvc := factloading.NewSvc(loadFactRepo, app.Logger, app.bathSize, txm, forecastStore)
	agentHndlr := agent.NewHndlr(app.Logger)
	r.Post("/agent", agentHndlr.AgentRoute)

	return r
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