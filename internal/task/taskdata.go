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
	getTaskQuery = "SELECT * FROM scheduler WHERE id = ?"

	getTasksQuery = "SELECT * FROM scheduler ORDER BY date LIMIT ?"

	getTasksByDateQuery = "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"

	getTasksBySearchStringQuery = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
	// Подготовленный запрос для обновления записи в базе данных.
	updateQuery = "UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?"
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

// getTasksByRows извлекает задачи из результата sql.Rows
func getTasksByRows(rows *sql.Rows) ([]Task, error) {
	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask получает задачу по ID
func (data TaskData) GetTask(id int) (Task, error) {

	row := data.db.QueryRow(getTaskQuery, id)

	var task Task
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return task, err
}

// GetTasks получает все задачи с ограничением по количеству
func (data TaskData) GetTasks(limit int) ([]Task, error) {

	rows, err := data.db.Query(getTasksQuery, limit)
	if err != nil {
		return nil, err
	}
	return getTasksByRows(rows)
}

// GetTasksByDate получает задачи по дате с ограничением по количеству
func (data TaskData) GetTasksByDate(date string, limit int) ([]Task, error) {

	rows, err := data.db.Query(getTasksByDateQuery, date, limit)
	if err != nil {
		return nil, err
	}
	return getTasksByRows(rows)
}

// GetTasksBySearchString получает задачи по поисковой строке с ограничением по количеству
func (data TaskData) GetTasksBySearchString(search string, limit int) ([]Task, error) {

	rows, err := data.db.Query(getTasksBySearchStringQuery, "%"+search+"%", "%"+search+"%", limit)
	if err != nil {
		return nil, err
	}
	return getTasksByRows(rows)
}

// UpdateTask обновляет задачу в базе данных.
func (data TaskData) UpdateTask(task Task) (bool, error) {

	// Начало транзакции.
	tx, err := data.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback() // Откат транзакции в случае ошибки.

	// Выполнение подготовленного запроса внутри транзакции.
	result, err := tx.Exec(updateQuery, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		return false, err
	}

	// Получение количества обновленных строк.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	// Коммит транзакции, если все операции без ошибок.
	if err = tx.Commit(); err != nil {
		return false, err
	}

	// Проверка, что была обновлена одна строка.
	return rowsAffected == 1, nil
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
