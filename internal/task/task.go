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

// List Структура представляет собой список задач
type List struct {
	Tasks []Task `json:"tasks"`
}

// Service представляет сервис для работы с задачами
type Service struct {
	taskData *TaskData
}

func sliceToTasks(list []Task) *List {
	if list == nil {
		return &List{Tasks: []Task{}}

	}
	return &List{Tasks: list}
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

// InitTaskService создает новый экземпляр Service
func InitTaskService(taskData *TaskData) Service {
	return Service{taskData: taskData}
}

// CreateTask Метод создает новую задачу
func (service Service) CreateTask(task Task) (int, error) {
	err := convertTask(&task)
	if err != nil {
		return 0, err
	}
	id, err := service.taskData.InsertTask(task)
	return int(id), err
}

func (service Service) UpdateTask(task Task) error {
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

func (service Service) GetTasks() (*List, error) {
	list, err := service.taskData.GetTasks(settings.TasksListRowsLimit)
	if err != nil {
		return nil, err
	}
	return sliceToTasks(list), err
}

func (service Service) SearchTasks(search string) (*List, error) {
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

func (service Service) GetTask(id string) (*Task, error) {
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

func (service Service) DeleteTask(id string) error {
	convId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	deleted, err := service.taskData.Delete(convId)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFoundTask
	}
	return nil
}

func (service Service) DoneTask(id string) error {
	convId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	task, err := service.taskData.GetTask(convId)
	if err != nil {
		return err
	}

	if len(task.Repeat) == 0 {
		deleted, err := service.taskData.Delete(convId)
		if err != nil {
			return err
		}
		if !deleted {
			return ErrNotFoundTask
		}
		return nil
	}

	task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
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
