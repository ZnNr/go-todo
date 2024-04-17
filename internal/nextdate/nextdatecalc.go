package nextdate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func CalculateNextDate(date, now time.Time, repeat string) (time.Time, error) {
	nextDate := date
	repeatParts := strings.Split(repeat, " ")

	switch repeatParts[0] {
	case "y":
		for {
			nextDate = nextDate.AddDate(1, 0, 0)
			if !nextDate.Before(now) {
				break
			}
		}

	case "d":
		if len(repeatParts) != 2 {
			return nextDate, errors.New("wrong repeat size")
		}
		days, err := strconv.Atoi(repeatParts[1])
		if err != nil {
			return nextDate, err
		}

		if days > 400 {
			return nextDate, errors.New("more than 400 days")
		}

		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if !nextDate.Before(now) {
				break
			}
		}

	default:
		return nextDate, errors.New("unsupported repeat type")
	}

	return nextDate, nil
}
