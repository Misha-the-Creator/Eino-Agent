package main

import (
	"fmt"
	"os"

	"github.com/Misha-the-Creator/Eino-Agent/internal/config"
	"github.com/Misha-the-Creator/Eino-Agent/internal/logger"
)

func main() {
	// ctx := context.Background()

	// Конфиг
	cfg := config.MustLoad()

	// Логгер
	log := logger.SetupLogger(cfg.App.IsLocal)
	log.Debug("Параметры конфигурации", "cfg", cfg)

	// Маршрутизатор HTTP
	appl := App{
		addr:        fmt.Sprintf(":%v", cfg.App.Port),
		
		Logger:      log,
	}

	if err := appl.run(appl.mount()); err != nil {
		log.Error("Произошла ошибка при запуске сервера:", "err", err)
		os.Exit(1)
	}
}