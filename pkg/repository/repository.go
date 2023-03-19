package repository

import (
	"github.com/foxleren/parser-bot/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Subscriber interface {
	CreateSubscriber(subscriber models.Subscriber) (int, error)
	GetAllSubscribers() ([]models.Subscriber, error)
	GetSubscriber(chatId int64) (models.Subscriber, error)
	DeleteSubscriber(chatId int64) error
}

type Repository struct {
	Subscriber
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Subscriber: NewSubscriberPostgres(db),
	}
}
