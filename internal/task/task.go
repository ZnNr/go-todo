package task

import (
	"errors"
	"github.com/ZnNr/go-todo/internal/nextdate"
	"github.com/ZnNr/go-todo/internal/settings"
	"strconv"
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

func sliceToTasks(list []Task) *TaskList {
	if list == nil {
		return &TaskList{Tasks: []Task{}}

	}
	return &TaskList{Tasks: list}
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

func (service TaskService) Update(task Task) error {
	err := convertTask(&task)
	if err != nil {
		return err
	}

	updated, err := service.taskData.UpdateTask(task)
	if err != nil {
		return err
	}
	if !updated {
		return ErrNotFoundTask
	}
	return nil
}

func (service TaskService) GetTasks() (*TaskList, error) {
	list, err := service.taskData.GetTasks(settings.TasksListRowsLimit)
	if err != nil {
		return nil, err
	}
	return sliceToTasks(list), err
}

func (service TaskService) SearchTasks(search string) (*TaskList, error) {
	date, err := time.Parse(settings.SearchDateFormat, search)
	if err == nil {
		list, err := service.taskData.GetTasksByDate(date.Format(settings.DateFormat), settings.TasksListRowsLimit)
		if err != nil {
			return nil, err
		}
		return sliceToTasks(list), nil
	}
	list, err := service.taskData.GetTasksBySearchString(search, settings.TasksListRowsLimit)
	return sliceToTasks(list), err
}

func (service TaskService) GetTask(id string) (*Task, error) {
	convId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	task, err := service.taskData.GetTask(convId)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
