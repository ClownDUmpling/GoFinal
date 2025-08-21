package api

import "errors"

// Создаем отдельный файл с константами для удобства работы с ними
const (
	DateFormat = "20060102"
)

type TaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskResponse struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

var (
	ErrInvalidRule        = errors.New("invalid repeat rule")
	ErrInvalidDate        = errors.New("invalid date format")
	ErrUnsupported        = errors.New("unsupported repeat rule")
	ErrInvalidDay         = errors.New("invalid day value")
	ErrInvalidMonth       = errors.New("invalid month value")
	ErrInvalidDayInterval = errors.New("invalid day interval: must be positive integer between 1-400")
)
