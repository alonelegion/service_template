package service_template

import (
	"context"
	"github.com/alonelegion/service_template/internal/app/config"
	"github.com/alonelegion/service_template/internal/shutdown_service"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"os"
	"time"

	application "github.com/alonelegion/service_template/internal/app"
	ll "github.com/alonelegion/service_template/internal/app/logger"
)

const (
	Name = "service-template"
)

func main() {
	// Генератор случайного числа с рандомизацией
	rand.Seed(time.Now().UnixNano())
	// Функция для выключения и перезагрузке сервисов
	ctx := shutdown_service.ShutDownContext(context.Background())

	// logger
	logger, err := ll.New(
		Name,
		os.Getenv("VERSION"),
		os.Getenv("ENV"),
		os.Getenv("LOG_LEVEL"),
	)

	if err != nil {
		// Ошибка при инициализации логгера
		log.Fatal("error while init logger", zap.Error(err))
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error("recover error", zap.Any("description", r))
		}
	}()

	// Получение конфигов
	appConfig, err := config.NewAppConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		logger.Fatal("error while init config", zap.Error(err))
	}

	app := application.New(
		Name,
		os.Getenv("VERSION"),
		os.Getenv("ENV"),
		appConfig,
		logger,
	)
}
