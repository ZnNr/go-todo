package task

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

const (
	driverName = "sqlite"

	tableSchema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY,
    date VARCHAR(8),
    title TEXT,
    comment TEXT,
    repeat VARCHAR(128)
);
`
	indexSchema = `
CREATE INDEX IF NOT EXISTS indexdate ON scheduler (date);
`
	insertQuery = `
INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)
`
)

// TaskData представляет структуру для работы с данными задач
type TaskData struct {
	db *sql.DB
}

// NewTaskData создает новый экземпляр TaskData с подключением к базе данных
func NewTaskData(dataSourceName string) (*TaskData, error) {
	db, err := openDb(dataSourceName)
	if err != nil {
		return nil, err
	}
	return &TaskData{db: db}, nil
}

// CloseDb закрывает соединение с базой данных
func (data *TaskData) CloseDb() {
	data.db.Close()
}

// InsertTask вставляет задачу в базу данных и возвращает ее ID
func (data *TaskData) InsertTask(task Task) (int64, error) {
	stmt, err := data.db.Prepare(insertQuery)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

// openDb открывает соединение с базой данных
func openDb(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(tableSchema); err != nil {
		return nil, err
	}
	if _, err := db.Exec(indexSchema); err != nil {
		return nil, err
	}
	return db, nil
}
