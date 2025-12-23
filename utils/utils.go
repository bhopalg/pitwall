package utils

import "time"

func ParseDate(date string) (*time.Time, error) {
	_parseDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}

	return &_parseDate, nil
}
