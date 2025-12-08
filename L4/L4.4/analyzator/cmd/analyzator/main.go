package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"analyzator/internal/app"
)

func main() {
	// Парсинг флагов командной строки
	port := flag.String("port", "8080", "Порт для запуска сервера (по умолчанию: 8080)")
	gcPercent := flag.String("gc", "-1", "Процент для GC (по умолчанию: -1, использовать настройки по умолчанию)")
	readTimeout := flag.Duration("read-timeout", 10*time.Second, "Таймаут чтения (по умолчанию: 10s)")
	writeTimeout := flag.Duration("write-timeout", 10*time.Second, "Таймаут записи (по умолчанию: 10s)")
	idleTimeout := flag.Duration("idle-timeout", 60*time.Second, "Таймаут простоя (по умолчанию: 60s)")

	flag.Parse()

	// Парсинг GC процента
	gcPercentInt, err := strconv.Atoi(*gcPercent)
	if err != nil {
		log.Fatalf("Ошибка парсинга GC процента: %v", err)
	}

	// Валидация GC процента
	if gcPercentInt < -1 {
		log.Fatalf("Недопустимое значение GC процента: %d (должно быть >= -1)", gcPercentInt)
	}

	// Создание конфигурации
	config := app.Config{
		Port:         *port,
		GCPercent:    gcPercentInt,
		ReadTimeout:  *readTimeout,
		WriteTimeout: *writeTimeout,
		IdleTimeout:  *idleTimeout,
	}

	// Вывод информации о запуске
	fmt.Println("=== Анализатор GC и памяти ===")
	fmt.Printf("Порт: %s\n", config.Port)
	if config.GCPercent >= 0 {
		fmt.Printf("GC процент: %d\n", config.GCPercent)
	} else {
		fmt.Printf("GC процент: по умолчанию\n")
	}
	fmt.Printf("Таймауты: чтение=%v, запись=%v, простой=%v\n",
		config.ReadTimeout, config.WriteTimeout, config.IdleTimeout)
	fmt.Println("==============================")

	// Создание и запуск приложения
	application := app.New(config)

	if err := application.Run(); err != nil {
		log.Fatalf("Ошибка запуска приложения: %v", err)
	}
}
