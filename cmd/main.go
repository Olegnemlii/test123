package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"githumb.com/olegnemlii/test123/internal/config"
	"githumb.com/olegnemlii/test123/internal/repository/postgres"
	"githumb.com/olegnemlii/test123/internal/service"
	"githumb.com/olegnemlii/test123/internal/transport/grpc/server"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключаем базу данных
	db, err := postgres.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Инициализируем репозиторий
	userRepo := postgres.NewUserRepository(db)

	// Создаем сервисы
	userService := service.NewUserService(userRepo)

	// Запускаем gRPC сервер
	grpcServer := server.NewGRPCServer(userService, cfg.GRPCPort)

	go func() {
		if err := grpcServer.Run(); err != nil {
			log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("Завершаем работу сервера...")
	grpcServer.Stop()
	log.Println("Сервер остановлен.")
}
