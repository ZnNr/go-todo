package task

import (
	"bytes"
	"encoding/json"
	"github.com/ZnNr/go-todo/internal/errorutil"
	"net/http"
)

var TaskServiceInstance Service

// taskFromRequestBody извлекает задачу из тела запроса
func taskFromRequestBody(r *http.Request) (Task, error) {
	var task Task

	buff := bytes.Buffer{}

	_, err := buff.ReadFrom(r.Body)
	if err != nil {
		return Task{}, err
	}

	err = json.Unmarshal(buff.Bytes(), &task)
	return task, err
}

func DonePostTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	err := TaskServiceInstance.DoneTask(id)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
	w.Write([]byte("{}"))
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	err := TaskServiceInstance.DeleteTask(id)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
	w.Write([]byte("{}"))
}

// PostTask обрабатывает POST запрос для создания задачи
func PostTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task, err := taskFromRequestBody(r)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}

	id, err := TaskServiceInstance.CreateTask(task)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}

	responseBody, err := json.Marshal(struct {
		Id int `json:"id"`
	}{Id: id})
	if err != nil {
		writeErrorAndRespond(w, http.StatusInternalServerError, err)
		return
	}

	_, err = w.Write(responseBody)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
}

// GetTask обрабатывает запрос на получение задачи по ID
func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Получаем значение параметра "id" из URL запроса
	id := r.URL.Query().Get("id")
	// Получаем задачу по ID с помощью сервиса TaskServiceInstance
	task, err := TaskServiceInstance.GetTask(id)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
	// Преобразуем полученную задачу в формат JSON и отправляем клиенту
	response, err := json.Marshal(task)
	if err != nil {
		writeErrorAndRespond(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(response)
}

// GetTasks обрабатывает запрос на получение списка задач
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Получаем значение параметра "search" из URL запроса
	search := r.URL.Query().Get("search")
	var tasks *List
	var err error
	// Если параметр "search" не указан, получаем все задачи, иначе ищем задачи по запросу
	if len(search) == 0 {
		tasks, err = TaskServiceInstance.GetTasks()
	} else {
		tasks, err = TaskServiceInstance.SearchTasks(search)
	}
	if err != nil {
		writeErrorAndRespond(w, http.StatusInternalServerError, err)
		return
	}
	// Преобразуем список задач в формат JSON и отправляем клиенту
	response, err := json.Marshal(tasks)
	if err != nil {
		writeErrorAndRespond(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(response)
}

func PutTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	task, err := taskFromRequestBody(r)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}

	err = TaskServiceInstance.UpdateTask(task)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
	w.Write([]byte("{}"))
}

// writeErrorAndRespond пишет ошибку в ответ и устанавливает соответствующий код состояния
func writeErrorAndRespond(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write(errorutil.MarshalError(err))
}
