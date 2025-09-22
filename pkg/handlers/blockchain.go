package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Index     int    `json:"index"`
	Timestamp string `json:"timestamp"`
	TouristID string `json:"tourist_id"`
	DataHash  string `json:"data_hash"`

	Itinerary json.RawMessage `json:"itinerary"`
	Emergency json.RawMessage `json:"emergency"`

	PrevHash string `json:"prev_hash"`
	Hash     string `json:"hash"`
}

var Blockchain []Block

func calculateHash(block Block) string {
	record := fmt.Sprint(block.Index) + block.Timestamp + block.TouristID +
		block.DataHash + string(block.Itinerary) + string(block.Emergency) + block.PrevHash
	h := sha256.Sum256([]byte(record))
	return hex.EncodeToString(h[:])
}

func generateBlock(prev Block, touristID, dataHash string, itinerary, emergency json.RawMessage) Block {
	block := Block{
		Index:     prev.Index + 1,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		TouristID: touristID,
		DataHash:  dataHash,
		Itinerary: itinerary,
		Emergency: emergency,
		PrevHash:  prev.Hash,
	}
	block.Hash = calculateHash(block)
	return block
}

func InitBlockchain() {
	genesis := Block{
		Index:     0,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Hash:      "GENESIS",
	}
	Blockchain = append(Blockchain, genesis)
}
