package database

import (
	"image"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

func LoadExistingFiles() {
	rows, _ := db.Query(`
		SELECT path FROM images
	`)
	defer rows.Close()

	for rows.Next() {
		var path string
		rows.Scan(&path)
		// You can process the loaded paths here
	}
}

func GetUndecidedImage() (string, error) {
	var path string
	err := db.QueryRow(`
		SELECT path FROM images WHERE decision = 'undecided' LIMIT 1
		`).Scan(&path)
	return path, err
}

func GetNextNImages(n int) ([]string, error) {
	var paths []string
	rows, err := db.Query("SELECT path FROM images WHERE decision IS NULL LIMIT ?", n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func GetUnchecksummedImagesFromDatabase() ([]string, error) {
	// var path string
	rows, err := db.Query(`
		SELECT path FROM images WHERE checksum is NULL
		`)
	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	defer rows.Close()
	return paths, err
}

func LoadSingleFile(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		logrus.Warn(err)
		// Handle error
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		logrus.Warn(err)
		// Handle error
	}
	return img
}

func GetUnchecksummedImagesFiles(imageCh chan<- image.Image) error {
	files, err := GetUnchecksummedImagesFromDatabase()
	if err != nil {
		return err
	}
	for _, file := range files {
		// Load a single file
		img := LoadSingleFile(file)
		imageCh <- img
	}
	return nil
}

func UpdateDecision(imagePath string, decision string) {
	statement, _ := db.Prepare(`
		UPDATE images SET decision = ? WHERE path = ?
		`)
	statement.Exec(decision, imagePath)
}

func CountNonchecksummedFiles() int64 {
	var count int64
	db.QueryRow(`
		SELECT COUNT(*) FROM images WHERE checksum is NULL
		`).Scan(&count)
	return count
}

func GetAllNonchecksummedFiles() int64 {
	var count int64
	db.QueryRow(`
		SELECT path FROM images WHERE checksum is NULL
	`).Scan(&count)
	return count
}

func UpdateChecksumInDatabase(imagePath string, checksum string) {
	statement, _ := db.Prepare(`
		UPDATE images SET checksum = ? WHERE path = ?
		`)
	statement.Exec(checksum, imagePath)
}

func InsertImagePathIntoDatabase(imagePath string) {
	// Insert the image path into the database
	insertSQL := `INSERT OR IGNORE INTO images (path) VALUES (?);`

	_, err := db.Exec(insertSQL, imagePath)
	if err != nil {
		log.Printf("Failed to insert image path %s: %v", imagePath, err)
	}
}

func InsertImagePathsIntoDatabase(imagePath []string) {
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

func GetImageBasePath() string {
	// Insert or update a setting
	return GetSetting("base_path")
}

func SetImageBasePath(basePath string) {
	// Insert or update a setting
	SetSetting("base_path", basePath)
}

func SetSetting(key string, value string) {
	// Insert or update a setting
	_, err := db.Exec(`INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)`, key, value)
	if err != nil {
		log.Fatal(err)
	}
}

func GetSetting(key string) string {
	// Insert or update a setting
	var value string
	err := db.QueryRow(`
		SELECT value FROM settings WHERE key = key LIMIT 1
		`).Scan(&value)
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func GetAllUndecidedPaths() ([]string, error) {
	// var path string
	rows, err := db.Query(`
		SELECT path FROM images WHERE decision is 'undecided'
		`)
	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	defer rows.Close()
	return paths, err
}
