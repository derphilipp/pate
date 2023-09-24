package database

import (
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var batchSize = 10

type FileChecksum struct {
	FilePath string
	Checksum string
}

func ChecksumWriter(checksumUpdates <-chan FileChecksum, progressCh chan<- int) {
	var count int = 0

	// Start the initial transaction
	tx, err := db.Begin()
	if err != nil {
		log.Panic("Error starting transaction:", err)
		return
	}
	// defer tx.Rollback()
	// Ensure any uncommitted changes are rolled back

	// Prepare the statement once outside the loop
	stmt, err := tx.Prepare("UPDATE images SET checksum = ? WHERE path = ?")
	if err != nil {
		log.Panic("Error preparing statement:", err)
		return
	}
	// defer stmt.Close() // Ensure the statement is closed after use

	for update := range checksumUpdates {
		_, err := stmt.Exec(update.Checksum, update.FilePath)
		if err != nil {
			log.Panic("Error executing statement:", err)
			fmt.Printf("OH GOD NO....")
			continue
		}

		count++
		if count >= batchSize {
			// ARGH
			if progressCh != nil {
				progressCh <- count
			}

			// Commit the current transaction
			if err := tx.Commit(); err != nil {
				log.Panic("Error committing transaction:", err)
				return
			}

			// Start a new transaction for the next batch
			tx, err = db.Begin()
			if err != nil {
				log.Panic("Error starting new transaction:", err)
				return
			}

			// Reuse the prepared statement for the new transaction
			stmt, err = tx.Prepare("UPDATE images SET checksum = ? WHERE path = ?")
			if err != nil {
				log.Panic("Error preparing statement for new transaction:", err)
				return
			}

			count = 0
		}
	}

	// Commit any remaining changes
	if count > 0 {
		if err := tx.Commit(); err != nil {
			log.Panic("Error committing final transaction:", err)
		}
	}
}
