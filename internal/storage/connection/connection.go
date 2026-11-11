package connection

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn struct {
	Host     string
	Port     string
	User     string
	Name     string
	Password string
	SSLMode  string
	Log      *slog.Logger
}

func (c *Conn) BuildDSN() string {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.PathEscape(c.User),
		url.PathEscape(c.Password),
		url.PathEscape(c.Host),
		url.PathEscape(c.Port),
		url.PathEscape(c.Name),
		url.PathEscape(c.SSLMode),
	)

	return dsn
}

func (c *Conn) MakeConnection(ctx context.Context) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, c.BuildDSN())
	if err != nil {
		c.Log.Error("Ошибка при создании соединения с БД", "err", err)
		os.Exit(1)
	}

	if err := pool.Ping(ctx); err != nil {
		c.Log.Error("Ошибка при попытке пропинговать БД", "err", err)
		os.Exit(1)
	}

	c.Log.Info("Пинговка БД прошла успешно")

	return pool
}