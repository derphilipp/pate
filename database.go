package main

import (
	"database/sql"

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

func updateDecision(imagePath string, decision string) {
	statement, _ := db.Prepare("UPDATE images SET decision = ? WHERE path = ?")
	statement.Exec(decision, imagePath)
}
