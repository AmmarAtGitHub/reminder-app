package models

import "time"

type Task struct {
	ID           int        `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	Description  *string    `json:"description" db:"description"`
	ReminderDate *time.Time `json:"reminder_date" db:"reminder_date"`
	Notified     bool       `json:"notified" db:"notified"`
	IsCompleted  bool       `json:"is_completed" db:"is_completed"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}
