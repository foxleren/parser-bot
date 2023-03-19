package models

type Subscriber struct {
	ID     int   `json:"id" db:"id"`
	ChatId int64 `json:"chat_id" db:"chat_id" binding:"required"`
}
