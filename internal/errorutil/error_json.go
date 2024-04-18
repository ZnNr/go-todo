package errorutil

import "encoding/json"

// MarshalError преобразует ошибку в формат JSON
func MarshalError(err error) []byte {
	type errJson struct {
		Error string `json:"error"`
	}
	res, _ := json.Marshal(errJson{Error: err.Error()})
	return res
}
