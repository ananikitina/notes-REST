package models

import "time"

type Note struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	UserID    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
