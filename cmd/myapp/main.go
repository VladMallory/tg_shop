package main

import (
	"log"
	"salle_parfume/internal/app"
)

func main() {
	// 1. сборка приложения
	myApp, err := app.New()
	if err != nil {
		log.Fatal("ошибка сборки приложения: %w", err)
	}

	// 1. запуск приложения
	myApp.Run()
}
