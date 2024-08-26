package models

import "time"

type Note struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	UserID    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
