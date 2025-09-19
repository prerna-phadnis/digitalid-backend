package main

import (
	"digitalid-backend/pkg/api"
	"fmt"
)

func main() {
	fmt.Println("Starting the server...")
	api.Start()
}
