package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/ping", ping())

	router.POST("/api/tourist/register", register())

	router.GET("/api/tourist/data/:id", BasicAuthMiddleware(), get())
}

var storage = make(map[string]RegisterRequest)

func ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	}
}

func register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data RegisterRequest

		// Parse and bind JSON to DTO
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Assign unique ID
		id := uuid.New().String()
		dataWithID := struct {
			ID string `json:"id"`
			RegisterRequest
		}{
			ID:              id,
			RegisterRequest: data,
		}

		// Marshal with ID included
		raw, err := json.MarshalIndent(dataWithID, "", "  ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
			return
		}

		// Generate SHA-256 hash of raw JSON
		hash := sha256.Sum256(raw)
		dataHash := hex.EncodeToString(hash[:])

		// Store in memory (temporary)
		storage[id] = data

		// Save JSON to local disk
		if err := os.WriteFile("./data/"+id+".json", raw, 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
			return
		}

		// TODO: push dataHash to blockchain
		fmt.Printf("hash %s\n", dataHash)

		// Generate QR Code containing the ID
		qr, err := qrcode.Encode(id, qrcode.Medium, 256)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR"})
			return
		}

		// Return QR code as PNG
		c.Header("Content-Type", "image/png")
		c.Writer.WriteHeader(http.StatusOK)
		_, _ = c.Writer.Write(qr)
	}
}

func get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if data, ok := storage[id]; ok {
			c.JSON(http.StatusOK, data)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
	}
}
