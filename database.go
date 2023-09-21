package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite", "./images.db")
	if err != nil {
		panic(err)
	}

	// Create table if not exists
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS images (path TEXT, checksum TEXT, decision TEXT DEFAULT 'undecided')")
	statement.Exec()

	statement, _ = db.Prepare("CREATE INDEX path_index ON images (path);")
	statement.Exec()

	statement, _ = db.Prepare("CREATE INDEX checksum_index ON images (checksum);")
	statement.Exec()
}

func loadExistingFiles() {
	rows, _ := db.Query("SELECT path FROM images")
	defer rows.Close()

	for rows.Next() {
		var path string
		rows.Scan(&path)
		// You can process the loaded paths here
	}
}

func getUndecidedImage() (string, error) {
	var path string
	err := db.QueryRow("SELECT path FROM images WHERE decision = 'undecided' LIMIT 1").Scan(&path)
	return path, err
}

func getUnchecksummedImagesFromDatabase() ([]string, error) {
	// var path string
	rows, err := db.Query("SELECT path FROM images WHERE checksum is NULL")
	defer rows.Close()
	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	return paths, err
}

func updateDecision(imagePath string, decision string) {
	statement, _ := db.Prepare("UPDATE images SET decision = ? WHERE path = ?")
	statement.Exec(decision, imagePath)
}

func countNonchecksummedFiles() int64 {
	var count int64
	db.QueryRow("SELECT COUNT(*) FROM images WHERE checksum is NULL").Scan(&count)
	return count
}

func getAllNonchecksummedFiles() int64 {
	var count int64
	db.QueryRow("SELECT path FROM images WHERE checksum is NULL").Scan(&count)
	return count
}

func updateChecksumInDatabase(imagePath string, checksum string) {
	statement, _ := db.Prepare("UPDATE images SET checksum = ? WHERE path = ?")
	statement.Exec(checksum, imagePath)
}

func insertImagePathIntoDatabase(imagePath string) {
	// Insert the image path into the database
	insertSQL := `INSERT OR IGNORE INTO images (path) VALUES (?);`

	_, err := db.Exec(insertSQL, imagePath)
	if err != nil {
		log.Printf("Failed to insert image path %s: %v", imagePath, err)
	}
}

func insertImagePathsIntoDatabase(imagePath []string) {
	// Insert the image path into the database
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	insertSQL := `INSERT OR IGNORE INTO images (path) VALUES (?);`
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		panic(err)
	}

	for _, path := range imagePath {
		_, err = stmt.Exec(path)
		if err != nil {
			panic(err)
		}
	}
	tx.Commit()
}
