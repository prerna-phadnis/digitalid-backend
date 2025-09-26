package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

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

		// Add computed fields (tourist name + digital expiry)
		dataWithID := struct {
			ID            string `json:"id"`
			TouristName   string `json:"tourist_name"`
			DigitalExpiry string `json:"digital_id_expiry"`
			RegisterRequest
		}{
			ID:              id,
			TouristName:     data.PersonalInfo.FullName,
			DigitalExpiry:   data.Travel.DepartureDate,
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

		// Save to PostgreSQL database
		query := `
		INSERT INTO tourist_data (id, tourist_name, digital_expiry, data_hash, raw_data)
		VALUES ($1, $2, $3, $4, $5)`

		_, err = db.Exec(query, id, data.PersonalInfo.FullName, data.Travel.DepartureDate, dataHash, raw)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database"})
			return
		}

		// Convert itinerary and emergency info into JSON for blockchain
		itineraryJSON, _ := json.Marshal(data.Travel.TripItinerary)
		emergencyJSON, _ := json.Marshal(data.Emergency.Contacts)

		// Add to blockchain
		if len(Blockchain) == 0 {
			InitBlockchain()
		}
		lastBlock := Blockchain[len(Blockchain)-1]
		newBlock := generateBlock(lastBlock, id, dataHash, itineraryJSON, emergencyJSON)
		Blockchain = append(Blockchain, newBlock)

		fmt.Printf("Tourist registered with ID: %s, Hash: %s\n", id, dataHash)

		// Prepare QR code payload
		qrPayload := struct {
			ID            string `json:"id"`
			TouristName   string `json:"tourist_name"`
			DigitalExpiry string `json:"digital_id_expiry"`
		}{
			ID:            id,
			TouristName:   data.PersonalInfo.FullName,
			DigitalExpiry: data.Travel.DepartureDate,
		}

		qrData, err := json.Marshal(qrPayload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare QR payload"})
			return
		}

		qr, err := qrcode.Encode(string(qrData), qrcode.Medium, 256)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR"})
			return
		}

		// Return QR code as PNG
		c.Header("Content-Type", "image/png")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"tourist_%s.png\"", id))
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(qr)
	}
}

func get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Query from PostgreSQL database
		query := `SELECT raw_data FROM tourist_data WHERE id = $1`
		var rawData []byte

		err := db.QueryRow(query, id).Scan(&rawData)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Tourist data not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data from database"})
			return
		}

		// Parse the JSON data
		var data interface{}
		if err := json.Unmarshal(rawData, &data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse stored data"})
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

// Additional endpoint to get QR code for existing tourist
func getQRCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Query tourist data
		query := `SELECT tourist_name, digital_expiry FROM tourist_data WHERE id = $1`
		var touristName, digitalExpiry string

		err := db.QueryRow(query, id).Scan(&touristName, &digitalExpiry)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Tourist not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tourist data"})
			return
		}

		// Generate QR code
		qrPayload := struct {
			ID            string `json:"id"`
			TouristName   string `json:"tourist_name"`
			DigitalExpiry string `json:"digital_id_expiry"`
		}{
			ID:            id,
			TouristName:   touristName,
			DigitalExpiry: digitalExpiry,
		}

		qrData, err := json.Marshal(qrPayload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare QR payload"})
			return
		}

		qr, err := qrcode.Encode(string(qrData), qrcode.Medium, 256)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
			return
		}

		c.Header("Content-Type", "image/png")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"tourist_%s.png\"", id))
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(qr)
	}
}

// Additional endpoint to view blockchain
func getBlockchain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"blockchain": Blockchain,
			"length":     len(Blockchain),
		})
	}
}
