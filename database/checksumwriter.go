package database

import (
	"log"

	_ "modernc.org/sqlite"
)

var batchSize = 10

type FileChecksum struct {
	FilePath string
	Checksum string
}

func ChecksumWriter(updates <-chan FileChecksum, progressCh chan<- int) {
	var count int

	// Start the initial transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return
	}
	defer tx.Rollback()
	// Ensure any uncommitted changes are rolled back

	// Prepare the statement once outside the loop
	stmt, err := tx.Prepare("UPDATE images SET checksum = ? WHERE path = ?")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return
	}
	defer stmt.Close() // Ensure the statement is closed after use

	for update := range updates {
		_, err := stmt.Exec(update.Checksum, update.FilePath)
		if err != nil {
			log.Println("Error executing statement:", err)
			continue
		}

		count++
		if count >= batchSize {
			progressCh <- count
			// Commit the current transaction
			if err := tx.Commit(); err != nil {
				log.Println("Error committing transaction:", err)
				return
			}

			// Start a new transaction for the next batch
			tx, err = db.Begin()
			if err != nil {
				log.Println("Error starting new transaction:", err)
				return
			}

			// Reuse the prepared statement for the new transaction
			stmt, err = tx.Prepare("UPDATE images SET checksum = ? WHERE path = ?")
			if err != nil {
				log.Println("Error preparing statement for new transaction:", err)
				return
			}

			count = 0
		}
	}

	// Commit any remaining changes
	if count > 0 {
		if err := tx.Commit(); err != nil {
			log.Println("Error committing final transaction:", err)
		}
	}
}
