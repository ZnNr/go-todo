package main

import (
	"github.com/ZnNr/go-todo/internal/nextdate"
	"github.com/ZnNr/go-todo/internal/settings"
	"github.com/ZnNr/go-todo/internal/task"
	"log"

	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	// Инициализация базы данных и задач.
	taskData, dbErr := task.NewTaskData(settings.Setting("TODO_DBFILE"))
	defer taskData.CloseDb()
	if dbErr != nil {
		log.Fatalf("Error initializing task data: %v", dbErr)
	}
	// Инициализация маршрутизатора.
	r := chi.NewRouter()
	task.TaskServiceInstance = task.InitTaskService(taskData)
	// Установка маршрутов для обработки файлов и API.
	r.Get("/*", FileServer)                      // Обработка запросов к файлам
	r.Get("/api/nextdate", nextdate.GetNextDate) // API для получения следующей даты
	r.Route("/api/task", func(r chi.Router) {
		r.Post("/", task.PostTask)   // API для создания задачи
		r.Get("/{id}", task.GetTask) // API для получения задачи по ID
	})
	r.Get("/api/tasks", task.GetTasks) // API для получения списка задач
	// Старт веб-сервера на указанном порту.
	port := settings.Setting("TODO_PORT")
	serverAddr := ":" + port
	log.Printf("Starting server on %s...", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Error starting the web server: %v", err)
	}
}
