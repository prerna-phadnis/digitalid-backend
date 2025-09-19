package api

import (
	"digitalid-backend/pkg/handlers"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Start initializes and runs the Gin server
func Start() {
	// db, err := database.InitDB()
	// if err != nil {
	// 	log.Fatal("Could not connect to the database: ", err)
	//}

	router := gin.Default()

	// r.LoadHTMLGlob("templates/*")

	handlers.SetupRoutes(router)
	fmt.Println("Server is running on port 8085")
	router.Run(":8085")
}
