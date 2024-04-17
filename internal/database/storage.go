package database

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

func InitializeDatabase(dataSource string) {
	log.Println("Conecting database...")
	dbFilePath := os.Getenv("TODO_DBFILE")
	if dbFilePath == "" {
		dbFilePath = "scheduler.db"
	}

	_, err := os.Stat(dbFilePath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	// Создание таблицы scheduler
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY,
            date TEXT,
            title TEXT,
            comment TEXT,
            repeat TEXT
        );
    `)
	if err != nil {
		log.Fatal(err)
	}

	// Создание индекса на колонку date
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created database schema")
}
