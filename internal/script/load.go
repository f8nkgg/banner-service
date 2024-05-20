package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection string
	connStr := "postgres://postgres:postgres@localhost:5433/banner?sslmode=disable"

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare the insert statement
	insertStmt := `
	INSERT INTO banners (tag_ids, feature_id, content, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Insert 1000 records
	for i := 1; i <= 1000; i++ {
		tagIDs := generateRandomTagIDs()
		featureID := i
		content := json.RawMessage(fmt.Sprintf(`{"data": "content for feature %d"}`, featureID))
		isActive := rand.Intn(2) == 1
		createdAt := time.Now()
		updatedAt := time.Now()

		_, err := db.Exec(insertStmt, pq.Array(tagIDs), featureID, content, isActive, createdAt, updatedAt)
		if err != nil {
			log.Fatalf("Error inserting record %d: %v", i, err)
		}
	}

	fmt.Println("Successfully inserted 1000 records.")
}

// generateRandomTagIDs generates a slice of random integers for tag_ids
func generateRandomTagIDs() []int {
	numTags := rand.Intn(10) + 1 // Random number of tags between 1 and 10
	tagIDs := make([]int, numTags)
	for i := 0; i < numTags; i++ {
		tagIDs[i] = rand.Intn(1000) + 1 // Random tag ID between 1 and 1000
	}
	return tagIDs
}
