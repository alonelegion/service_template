package service_template

import (
	"context"
	"github.com/alonelegion/service_template/internal/shutdown_service"
	"math/rand"
	"time"
)

func main() {
	// Генератор случайного числа с рандомизацией
	rand.Seed(time.Now().UnixNano())
	// Функция для выключения и перезагрузке сервисов
	ctx := shutdown_service.ShutDownContext(context.Background())

	// logger
}
