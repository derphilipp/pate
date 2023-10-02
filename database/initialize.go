package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDatabase() {
	var err error
	db, err = sql.Open("sqlite", "./images.db")
	if err != nil {
		panic(err)
	}

	// Create table if not exists
	statement, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS images
			(
				path TEXT PRIMARY KEY,
				checksum TEXT,
				decision TEXT DEFAULT 'undecided',
				duplicate TEXT DEFAULT 'no',
				valid TEXT DEFAULT 'unknown'
			)
		`)
	statement.Exec()

	statement, _ = db.Prepare(`
		CREATE TABLE IF NOT EXISTS settings
			(key TEXT PRIMARY KEY, value TEXT)
		`)
	statement.Exec()

	statement, _ = db.Prepare(`
		CREATE INDEX path_index ON images (path);
		`)
	statement.Exec()

	statement, _ = db.Prepare(
		`CREATE INDEX checksum_index ON images (checksum)
		`)
	statement.Exec()
}
