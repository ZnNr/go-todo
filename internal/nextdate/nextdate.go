// Package nextdate реализует функцию NextDate для вычисления следующей даты
// на основе заданной даты и правила повторения. Поддерживаются различные
// типы правил, такие как "d", "y", "w", "m".
package nextdate

import (
	"errors"
	"github.com/ZnNr/go-todo/internal/settings"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotFoundRule = errors.New("Not found repeat rule") // ErrNotFoundRule возвращается, когда не удалось найти правило повторения.
	ErrBadRule      = errors.New("Bad repeat rule")       // ErrBadRule возвращается в случае некорректного правила повторения.
)

// rules содержит маппинг регулярных выражений на функции, обрабатывающие правила повторения.
var rules = map[*regexp.Regexp]func(now time.Time, date, repeat string) (string, error){
	regexp.MustCompile("^d\\s\\d{1,3}$"):                                     dayRule,
	regexp.MustCompile("^y$"):                                                yearRule,
	regexp.MustCompile("^w\\s[1-7]?(,[1-7]){0,6}$"):                          weekRule,
	regexp.MustCompile("^m\\s-?\\d+(,-?\\d+){0,30}(\\s\\d+(,\\d+){0,11})?$"): monthRule,
}

// months преобразует строки с номерами месяцев в булев массив, основанный на их наличии.
func months(months []string) ([12]bool, error) {
	ans := [12]bool{}
	for _, month := range months {
		monthIndex, err := strconv.Atoi(month)
		if err != nil || monthIndex > 12 || monthIndex < 1 {
			return ans, ErrBadRule
		}
		ans[monthIndex-1] = true
	}
	return ans, nil
}

// inSet проверяет, присутствует ли текущий день во множестве дней.
func inSet(now time.Time, daysSet map[int]bool) bool {
	if daysSet[now.Day()] {
		return true
	}
	year, month, _ := now.Date()
	firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	return daysSet[-1] && firstDayOfMonth.AddDate(0, 1, -1).Equal(now) ||
		daysSet[-2] && firstDayOfMonth.AddDate(0, 1, -2).Equal(now)
}

// Функция monthRule обрабатывает правило повторения для месяца.
func monthRule(now time.Time, date, repeat string) (string, error) {
	splitRepeat := strings.Split(repeat, " ")
	monthDays := strings.Split(splitRepeat[1], ",")
	includeDaysSet := map[int]bool{}
	for _, monthDay := range monthDays {
		dayIndex, _ := strconv.Atoi(monthDay)
		if dayIndex < -2 || dayIndex > 31 {
			return "", ErrBadRule
		}
		includeDaysSet[dayIndex] = true
	}

	var includeMonths [12]bool
	if len(splitRepeat) > 2 {
		var err error
		includeMonths, err = months(strings.Split(splitRepeat[2], ","))
		if err != nil {
			return "", err
		}
	} else {
		includeMonths = [12]bool{}
		for i := range includeMonths {
			includeMonths[i] = true
		}
	}
	curDate, err := time.Parse(settings.DateFormat, date)
	if err != nil {
		return "", err
	}

	if now.After(curDate) {
		curDate = now
	}

	curDate = curDate.AddDate(0, 0, 1)
	for !includeMonths[curDate.Month()-1] || !inSet(curDate, includeDaysSet) {
		curDate = curDate.AddDate(0, 0, 1)
	}
	return curDate.Format(settings.DateFormat), nil
}

// Функция weekRule обрабатывает правило повторения для недели.
func weekRule(now time.Time, date, repeat string) (string, error) {
	weekDays := strings.Split(strings.Split(repeat, " ")[1], ",")
	weekDaysSet := map[int]bool{}
	for _, weekDay := range weekDays {
		dayIndex, _ := strconv.Atoi(weekDay)
		weekDaysSet[dayIndex%7] = true
	}

	curDate, err := time.Parse(settings.DateFormat, date)
	if err != nil {
		return "", err
	}

	if now.After(curDate) {
		curDate = now
	}

	curDate = curDate.AddDate(0, 0, 1)

	for !weekDaysSet[int(curDate.Weekday())] {
		curDate = curDate.AddDate(0, 0, 1)
	}

	return curDate.Format(settings.DateFormat), nil
}

// Функция dayRule обрабатывает правило повторения для дня.
func dayRule(now time.Time, date, repeat string) (string, error) {
	items := strings.Split(repeat, " ")

	days, err := strconv.Atoi(items[1])
	if err != nil {
		return "", err
	}

	if days > 400 || days < 1 {
		return "", ErrBadRule
	}

	curDate, err := time.Parse(settings.DateFormat, date)
	if err != nil {
		return "", err
	}

	for next := true; next; next = !curDate.After(now) {
		curDate = curDate.AddDate(0, 0, days)
	}

	return curDate.Format(settings.DateFormat), nil
}

// Функция yearRule обрабатывает правило повторения для дня.
func yearRule(now time.Time, date, _ string) (string, error) {

	curDate, err := time.Parse(settings.DateFormat, date)
	if err != nil {
		return "", err
	}

	for next := true; next; next = !curDate.After(now) {
		curDate = curDate.AddDate(1, 0, 0)
	}

	return curDate.Format(settings.DateFormat), nil
}

// NextDate принимает текущее время `now`, строку `date` и строку `repeat`.
// Если `repeat` не является пустой строкой, функция ищет соответствующее правило повторения в `rules`,
// затем выполняет соответствующую функцию для обработки правила повторения и возвращает результат.
// Если правило не найдено, возвращается ошибка `ErrNotFoundRule`.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", nil
	}

	for pattern, f := range rules {
		if pattern.MatchString(repeat) {
			result, err := f(now, date, repeat)
			if err != nil {
				return "", err
			}
			return result, nil
		}
	}

	return "", ErrNotFoundRule
}
