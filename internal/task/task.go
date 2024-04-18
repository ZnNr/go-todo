package task

import (
	"errors"
	"github.com/ZnNr/go-todo/internal/nextdate"
	"github.com/ZnNr/go-todo/internal/settings"
	"time"
)

var (
	ErrRequireTitle = errors.New("require task title")
	ErrNotFoundTask = errors.New("not found task")
)

// Task Структура представляет собой модель задачи
type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// TaskList Структура представляет собой список задач
type TaskList struct {
	Tasks []Task `json:"tasks"`
}

// TaskService представляет сервис для работы с задачами
type TaskService struct {
	taskData *TaskData
}

// Функция convertTask конвертирует и проверяет задачу перед сохранением
func convertTask(task *Task) error {
	if len(task.Title) == 0 {
		return ErrRequireTitle
	}
	// Установка даты по умолчанию, если она не была указана, и проверка формата даты
	now := time.Now().Format(settings.DateFormat)
	if len(task.Date) == 0 {
		task.Date = now
	}
	_, err := time.Parse(settings.DateFormat, task.Date)
	if err != nil {
		return err
	}
	// Рассчет и установка следующей даты, если необходимо
	nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	if task.Date < now {
		if len(nextDate) == 0 {
			task.Date = now
		} else {
			task.Date = nextDate
		}
	}
	return nil
}

// InitTaskService создает новый экземпляр TaskService
func InitTaskService(taskData *TaskData) TaskService {
	return TaskService{taskData: taskData}
}

// CreateTask Метод создает новую задачу
func (service TaskService) CreateTask(task Task) (int, error) {
	err := convertTask(&task)
	if err != nil {
		return 0, err
	}
	id, err := service.taskData.InsertTask(task)
	return int(id), err
}
