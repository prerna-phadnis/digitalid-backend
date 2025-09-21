package api

import (
	"digitalid-backend/pkg/handlers"
	"fmt"

	"github.com/gin-gonic/gin"
)


func Start() {
	
	handlers.InitBlockchain()

	router := gin.Default()

	
	handlers.SetupRoutes(router)

	fmt.Println("Server is running on port 8085")
	router.Run(":8085")
}
