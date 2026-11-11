package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Misha-the-Creator/Eino-Agent/internal/config"
	"github.com/Misha-the-Creator/Eino-Agent/internal/logger"
	"github.com/Misha-the-Creator/Eino-Agent/internal/storage/connection"
)

func main() {
	ctx := context.Background()

	// Конфиг
	cfg := config.MustLoad()

	// Логгер
	log := logger.SetupLogger(cfg.App.IsLocal)
	log.Debug("Параметры конфигурации", "cfg", cfg)

	// БД
	dbConnection := connection.Conn{
	Host:     cfg.DB.Host,
	Port:     cfg.DB.Port,
	User:     cfg.DB.User,
	Name:     cfg.DB.Name,
	Password: cfg.DB.Password,
	SSLMode:  cfg.DB.SSLMode,
	Log:      log,
	}
	pool := dbConnection.MakeConnection(ctx)
	defer pool.Close()

	// Маршрутизатор HTTP
	appl := App{
		ctx:       ctx,
		addr:      fmt.Sprintf(":%v", cfg.App.Port),
		apiKey:    cfg.LLM.APIkey,
		modelName: cfg.LLM.ModelName,
		Logger:    log,
		pool: pool,
	}

	h, err := appl.mount()
	if err != nil {
		log.Error("Произошла ошибка при сборке роутера:", "err", err)
		os.Exit(1)
	}

	if err := appl.run(h); err != nil {
		log.Error("Произошла ошибка при запуске сервера:", "err", err)
		os.Exit(1)
	}
}
