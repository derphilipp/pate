package main

import (
	"log"

	_ "modernc.org/sqlite"
)

type FileChecksum struct {
	FilePath string
	Checksum string
}

func dbWriter(updates <-chan FileChecksum) {
	var count int

	// Start the initial transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return
	}
	defer tx.Rollback() // Ensure any uncommitted changes are rolled back

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

/*
func dbWriter(updates <-chan FileChecksum) {
	var count int
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("UPDATE images SET checksum = ? WHERE path = ?")
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		_, err := stmt.Exec(update.Checksum, update.FilePath)
		if err != nil {
			log.Fatal(err)
		}

		count++
		if count >= batchSize {
			if err := tx.Commit(); err != nil {
				log.Fatal(err)
			}
			stmt.Close()
			tx, err = db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			stmt, err = tx.Prepare("UPDATE images SET checksum = ? WHERE path = ?")
			if err != nil {
				log.Fatal(err)
			}

			count = 0
		}
	}

	if count > 0 {
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
	}
}
*/

// func dbWriter(batchChan <-chan FileChecksum) {
// 	for {
// 		tx, err := db.Begin()
// 		if err != nil {
// 			fmt.Println("Error starting transaction:", err)
// 			continue
// 		}

// 		stmt, err := tx.Prepare("UPDATE images SET checksum = ? WHERE path = ?")
// 		if err != nil {
// 			fmt.Println("Error preparing statement:", err)
// 			continue
// 		}

// 		// read from channel
// 		fc := <-batchChan

// 		_, err := stmt.Exec(fc.FilePath, fc.Checksum)
// 		if err != nil {
// 			fmt.Println("Error inserting into database:", err)
// 			break
// 		}

// 				stmt.Close()
// 				c := tx.Commit()
// 				if c != nil {
// 					fmt.Println("Error committing transaction:", c)
// 				}
// 			}
// 		}
// 	}
// }
