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
	taskData, err := task.NewTaskData(settings.Setting("TODO_DBFILE"))
	defer taskData.CloseDb()
	if err != nil {
		panic(err)
	}
	// Инициализация маршрутизатора.
	r := chi.NewRouter()
	// Маршруты для обработки файлов и API.
	r.Get("/*", FileServer)
	r.Get("/api/nextdate", nextdate.GetNextDate)
	task.TaskServiceInstance = task.InitTaskService(taskData)
	r.Post("/api/task", task.PostTask)
	// Старт веб-сервера на указанном порту.
	port := settings.Setting("TODO_PORT")
	log.Printf("Server started on port :%s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Error starting the web server: %v", err)
	}
}
