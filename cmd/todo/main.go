package main

import (
	"github.com/ZnNr/go-todo/internal/authorization"
	"github.com/ZnNr/go-todo/internal/nextdate"
	"github.com/ZnNr/go-todo/internal/settings"
	"github.com/ZnNr/go-todo/internal/task"
	"log"

	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	// Инициализация базы данных и задач.
	dbFile := settings.Setting("TODO_DBFILE")
	taskData, dbErr := task.NewTaskData(dbFile)
	defer taskData.CloseDb()
	if dbErr != nil {
		log.Fatalf("Error initializing task data: %v", dbErr)
	}

	// Инициализация маршрутизатора.
	r := chi.NewRouter()

	// Инициализация службы задач.
	task.TaskServiceInstance = task.InitTaskService(taskData)

	// Установка маршрутов для обработки файлов и API.
	r.Get("/*", FileServer) // Обработка запросов к файлам

	pass := settings.Setting("TODO_PASSWORD")
	secretKey := settings.Setting("SECRET_KEY")

	// Инициализация службы авторизации.
	authorization.Service = authorization.InitSignService(pass, []byte(secretKey))

	// Регистрация маршрута API для аутентификации пользователя.
	r.Post("/api/signin", authorization.PostPass)

	r.Get("/api/nextdate", nextdate.GetNextDate) // API для получения следующей даты

	// Группировка маршрутов для задач с общей авторизацией.
	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if len(pass) == 0 {
					next.ServeHTTP(w, req)
				} else {
					authMiddleware := &authorization.AuthMiddleware{} // Создание экземпляра AuthMiddleware
					authMiddleware.Auth(next).ServeHTTP(w, req)       // Использование AuthMiddleware
				}
			})
		})

		r.Post("/api/task", task.PostTask)          // Создание задачи
		r.Put("/api/task", task.PutTask)            // Обновление задачи
		r.Delete("/api/task", task.DeleteTask)      // Удаление задачи
		r.Get("/api/task", task.GetTask)            // Получение конкретной задачи
		r.Post("/api/task/done", task.DonePostTask) // Отметка задачи как выполненной
		r.Get("/api/tasks", task.GetTasks)          // API для получения списка задач
	})

	// Старт веб-сервера на указанном порту.
	port := settings.Setting("TODO_PORT")
	serverAddr := ":" + port
	log.Printf("Starting server on %s...", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Error starting the web server: %v", err)
	}
}
