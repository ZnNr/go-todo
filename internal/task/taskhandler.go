package task

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var TaskServiceInstance TaskService

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
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}

	_, err = w.Write(responseBody)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}
}

// writeErrorAndRespond пишет ошибку в ответ и устанавливает соответствующий код состояния
func writeErrorAndRespond(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write(MarshalError(err))
}

// MarshalError преобразует ошибку в формат JSON
func MarshalError(err error) []byte {
	type errJson struct {
		Error string `json:"error"`
	}
	res, _ := json.Marshal(errJson{Error: err.Error()})
	return res
}
