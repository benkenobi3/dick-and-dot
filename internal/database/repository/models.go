package repository

type Dick struct {
	UserID    int64 `json:"userId"`
	ChatID    int64 `json:"chatId"`
	Length    int64 `json:"length"`
	UpdatedAt int64 `json:"updatedAt"`
}
