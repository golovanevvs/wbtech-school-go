package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	pkgPrometheus "analyzator/internal/pkg/pkgPrometheus"
	"analyzator/internal/transport"
)

// Config конфигурация приложения
type Config struct {
	Port         string
	GCPercent    int // Процент для GC, -1 для отключения принудительного GC
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// App структура приложения
type App struct {
	server *http.Server
	config Config
}

// New создает новое приложение
func New(config Config) *App {
	return &App{
		config: config,
	}
}

// Setup настраивает приложение
func (a *App) Setup() error {
	// Настройка GC
	if a.config.GCPercent >= 0 {
		prev := debug.SetGCPercent(a.config.GCPercent)
		log.Printf("Установлен GC процент: %d (предыдущий: %d)", a.config.GCPercent, prev)
	}

	// Инициализация Prometheus метрик
	pkgPrometheus.Init()
	log.Println("Prometheus метрики инициализированы")

	// Создание HTTP сервера
	router := transport.NewRouter()

	a.server = &http.Server{
		Addr:         ":" + a.config.Port,
		Handler:      router,
		ReadTimeout:  a.config.ReadTimeout,
		WriteTimeout: a.config.WriteTimeout,
		IdleTimeout:  a.config.IdleTimeout,
	}

	return nil
}

// Start запускает сервер
func (a *App) Start() error {
	log.Printf("Запуск сервера на порту %s", a.config.Port)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	return nil
}

// Stop останавливает сервер
func (a *App) Stop(ctx context.Context) error {
	log.Println("Остановка сервера...")
	return a.server.Shutdown(ctx)
}

// Run запускает приложение с обработкой сигналов
func (a *App) Run() error {
	// Настройка приложения
	if err := a.Setup(); err != nil {
		return fmt.Errorf("ошибка настройки приложения: %w", err)
	}

	// Запуск сервера
	if err := a.Start(); err != nil {
		return fmt.Errorf("ошибка запуска сервера: %w", err)
	}

	// Ожидание сигналов для корректного завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Получен сигнал завершения...")

	// Корректное завершение
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.Stop(ctx); err != nil {
		return fmt.Errorf("ошибка остановки сервера: %w", err)
	}

	log.Println("Сервер остановлен")
	return nil
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		Port:         "8080",
		GCPercent:    -1, // -1 означает использовать настройки по умолчанию
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
