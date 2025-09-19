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

	router.GET("/api/tourist/data/:id", get())
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
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		data.ID = uuid.New().String()
		raw, _ := json.Marshal(data)
		hash := sha256.Sum256(raw)
		dataHash := hex.EncodeToString(hash[:])

		// store locally (TODO: later in s3/ipfs)
		storage[data.ID] = data

		if err := os.WriteFile("./data/"+data.ID+".json", raw, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
			return
		}

		// TODO: push datahash to blockchain

		fmt.Printf("hash %s\n", dataHash)

		// Generate QR Code PNG
		qr, _ := qrcode.Encode(data.ID, qrcode.Medium, 256)

		c.Header("Content-Type", "image/png")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(qr)
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
