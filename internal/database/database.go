package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func InitializeDatabase() {
	dbFilePath := os.Getenv("TODO_DBFILE")
	if dbFilePath == "" {
		dbFilePath = "scheduler.db"
	}

	_, err := os.Stat(dbFilePath)
	if os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", dbFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		_, err = db.Exec(`
            CREATE TABLE IF NOT EXISTS scheduler (
                id INTEGER PRIMARY KEY,
                date TEXT,
                title TEXT,
                comment TEXT,
                repeat TEXT
            );
            CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
        `)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Created database schema")
	} else if err != nil {
		log.Fatal(err)
	}

}
