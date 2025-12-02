package main

import (
	"os"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/app"
)

func main() {
	// Создаём экземпляр приложения
	application := app.NewApp()

	// Запускаем приложение
	if err := application.Run(); err != nil {
		os.Exit(1)
	}
}
