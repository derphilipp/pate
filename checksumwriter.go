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
