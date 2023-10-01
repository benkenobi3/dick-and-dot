package repository

import "time"

type Dick struct {
	UserID    int64     `json:"userId"`
	ChatID    int64     `json:"chatId"`
	Length    int64     `json:"length"`
	UpdatedAt time.Time `json:"updatedAt"`
}
