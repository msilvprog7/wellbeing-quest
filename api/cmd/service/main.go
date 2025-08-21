package main

import (
	"api.wellbeingquest.app/internal/services"
)

func main() {
	// Resource setup
	services.LoadDotEnv(".env")
	dataHandler := services.GetDataHandler()

	// Gin setup
	router := services.SetupRouter(dataHandler)
	router.Run("localhost:8080")
}
