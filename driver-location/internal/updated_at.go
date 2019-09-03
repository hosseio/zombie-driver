package domain

import (
	"errors"
	"time"
)

type UpdatedAt struct {
	date time.Time
}

var ErrUpdatedAtCanNotBeFuture = errors.New("the date is future")
var ErrUpdatedAtInvalidFormat = errors.New("the format should be RFC3339")

func FromString(date string) (UpdatedAt, error) {
	timeDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return UpdatedAt{}, ErrUpdatedAtInvalidFormat
	}

	return NewUpdatedAt(timeDate)
}

func NewUpdatedAt(timeDate time.Time) (UpdatedAt, error) {
	updatedAt := UpdatedAt{timeDate}
	err := updatedAt.validate()
	if err != nil {
		return UpdatedAt{}, err
	}

	return updatedAt, nil
}

func (u UpdatedAt) validate() error {
	if u.date.After(time.Now()) {
		return ErrUpdatedAtCanNotBeFuture
	}

	return nil
}

func (u UpdatedAt) Date() time.Time {
	return u.date
}

func (u UpdatedAt) String() string {
	return u.Date().Format(time.RFC3339)
}
