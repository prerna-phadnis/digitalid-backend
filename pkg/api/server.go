package api

import (
	"digitalid-backend/pkg/handlers"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func Start() {
	dsn := "host=localhost port=5432 user=postgres password=password dbname=tourist_db sslmode=disable"

	if err := handlers.InitDB(dsn); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer handlers.CloseDB()

	handlers.InitBlockchain()

	router := gin.Default()

	handlers.SetupRoutes(router)

	fmt.Println("Server is running on port 8085")
	router.Run(":8085")
}
