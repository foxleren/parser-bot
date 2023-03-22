package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errUnableToSubscribe   = errors.New("error: failed while subscribing")
	errUnableToUnsubscribe = errors.New("error: failed while unsubscribing")
)

const (
	unableToSubscribeText   = "Не удалось подписаться на рассылку либо вы уже подписаны. Для проверки подписки используйте команду \n/check_subscription."
	unableToUnsubscribeText = "Не удалось отписаться от рассылки либо вы уже отписаны. Для проверки подписки используйте команду \n/check_subscription."
	unknownError            = "Произошла неизвестная ошибка."
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, "")
	switch err {
	case errUnableToSubscribe:
		msg.Text = unableToSubscribeText
	case errUnableToUnsubscribe:
		msg.Text = unableToUnsubscribeText
	default:
		msg.Text = unknownError
	}
	b.bot.Send(msg)
}
