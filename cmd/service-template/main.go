package service_template

import (
	"context"
	"github.com/alonelegion/service_template/internal/shutdown_service"
	"log"
	"math/rand"
	"os"
	"time"

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
		log.Fatal("error while init logger")
	}
}
