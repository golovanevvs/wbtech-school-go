package main

import (
	"os"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/app"
)

func main() {
	application := app.New()

	if err := application.Run(); err != nil {
		os.Exit(1)
	}
}
