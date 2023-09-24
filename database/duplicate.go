package database

import (
	"fmt"
	"log"
)

func DetectAndHandleDuplicates() {
	rows, err := db.Query(`SELECT path, checksum FROM images WHERE checksum IS NOT NULL`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	checksumMap := make(map[string][]string)

	for rows.Next() {
		var path, checksum string
		if err := rows.Scan(&path, &checksum); err != nil {
			log.Fatal(err)
		}
		checksumMap[checksum] = append(checksumMap[checksum], path)
	}

	for _, paths := range checksumMap {
		if len(paths) > 1 {
			// Mark all but one file as duplicate
			for i := 1; i < len(paths); i++ {
				_, err := db.Exec(`UPDATE images SET duplicate = true WHERE path = ?`, paths[i])
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Marked %s as duplicate\n", paths[i])
			}
		}
	}
}
