package repository

import (
	"fmt"
	"github.com/foxleren/parser-bot/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/siruspen/logrus"
)

type SubscriberPostgres struct {
	db *sqlx.DB
}

func NewSubscriberPostgres(db *sqlx.DB) *SubscriberPostgres {
	return &SubscriberPostgres{db: db}
}

func (s *SubscriberPostgres) CreateSubscriber(subscriber models.Subscriber) (int, error) {
	var id int
	createSubscriberQuery := fmt.Sprintf("INSERT INTO %s (chat_id) VALUES ($1) RETURNING id", subscribersTable)
	row := s.db.QueryRow(createSubscriberQuery, subscriber.ChatId)

	if err := row.Scan(&id); err != nil {
		logrus.Printf("repo: CreateSubscriber(): %v", err.Error())
		return 0, err
	}

	return id, nil
}

func (s *SubscriberPostgres) GetAllSubscribers() ([]models.Subscriber, error) {
	logrus.Printf("%v", s.db != nil)

	var subscribers []models.Subscriber
	getAllQuery := fmt.Sprintf("SELECT id, chat_id FROM %s", subscribersTable)
	err := s.db.Select(&subscribers, getAllQuery)

	if err != nil {
		logrus.Printf("repo: GetAllSubscribers(): %v", err.Error())
	}

	return subscribers, err
}

func (s *SubscriberPostgres) GetSubscriber(chatId int64) (models.Subscriber, error) {
	var subscriber models.Subscriber
	query := fmt.Sprintf("SELECT id, chat_id FROM %s WHERE chat_id = $1", subscribersTable)
	err := s.db.Get(&subscriber, query, chatId)

	if err != nil {
		logrus.Printf("repo: GetSubscriber(): %v", err.Error())
	}

	return subscriber, err
}

func (s *SubscriberPostgres) DeleteSubscriber(chatId int64) error {
	deleteCartItemByIDQuery := fmt.Sprintf("DELETE FROM %s WHERE chat_id = %d", subscribersTable, chatId)
	_, err := s.db.Exec(deleteCartItemByIDQuery)

	if err != nil {
		logrus.Printf("repo: DeleteSubscriber(): %v", err.Error())
	}

	return err
}
