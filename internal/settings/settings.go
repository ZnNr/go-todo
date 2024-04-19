// Package settings предоставляет функциональность для работы с настройками приложения.
package settings

import "os"

// DateFormat представляет формат даты по умолчанию.
var DateFormat = "20060102"

var SearchDateFormat = "02.01.2006"

var TasksListRowsLimit = 50

// defaultEnv содержит значения по умолчанию для некоторых настроек.
var defaultEnv = map[string]string{
	"TODO_PORT":     "7540",
	"TODO_DBFILE":   "../scheduler.db",
	"TODO_PASSWORD": "",
	"SECRET_KEY":    "my_secret_key",
}

// Setting возвращает значение настройки для указанного ключа.
// Если значение не задано в переменных окружения, то используется значение по умолчанию.
func Setting(key string) string {
	value := os.Getenv(key)
	if len(value) > 0 {
		return value
	}
	return defaultEnv[key]
}

// WebPath содержит путь к директории с статическими файлами для веб-сервера.
const WebPath = "./web/"
