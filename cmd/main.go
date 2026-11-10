package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Misha-the-Creator/Eino-Agent/internal/config"
	"github.com/Misha-the-Creator/Eino-Agent/internal/logger"
)

func main() {
	ctx := context.Background()

	// Конфиг
	cfg := config.MustLoad()

	// Логгер
	log := logger.SetupLogger(cfg.App.IsLocal)
	log.Debug("Параметры конфигурации", "cfg", cfg)

	// Маршрутизатор HTTP
	appl := App{
		ctx:       ctx,
		addr:      fmt.Sprintf(":%v", cfg.App.Port),
		apiKey:    cfg.LLM.APIkey,
		modelName: cfg.LLM.ModelName,
		Logger:    log,
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
