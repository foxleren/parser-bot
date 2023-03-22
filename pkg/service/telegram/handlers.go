package telegram

import (
	"github.com/foxleren/parser-bot/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/siruspen/logrus"
)

const (
	commandStart          = "start"
	commandSubscribe      = "subscribe"
	commandCheckSubscribe = "check_subscription"
	commandUnsubscribe    = "unsubscribe"

	replyStart          = "Добро пожаловать!\nДля подписки на рассылку введи команду /subscribe. Далее вы будете ежедневно получать новые данные."
	replyUnknownCommand = "Неизвестная команда."

	successfulSubscription   = "Вы успешно подписались на рассылку!"
	successfulUnsubscription = "Вы успешно отписались от рассылки."
	subscriptionStatusGood   = "Статус подписки: активирована."
	subscriptionStatusBad    = "Статус подписки: деактивирована."
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleCommandStart(message)
	case commandSubscribe:
		return b.handleCommandSubscribe(message)
	case commandUnsubscribe:
		return b.handleCommandUnsubscribe(message)
	case commandCheckSubscribe:
		return b.handleCommandCheckSubscribe(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

//func (b *Bot) handleMessage(message *tgbotapi.Message) error {
//	return nil
//}

func (b *Bot) handleCommandStart(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, replyStart)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, replyUnknownCommand)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleCommandSubscribe(message *tgbotapi.Message) error {
	subscriber := models.Subscriber{ChatId: message.Chat.ID}

	var id int
	id, err := b.repo.CreateSubscriber(subscriber)
	if err != nil {
		return errUnableToSubscribe
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, successfulSubscription)
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	logrus.Printf("Created id: %d", id)

	b.sendData(message.Chat.ID)

	return nil
}

func (b *Bot) handleCommandUnsubscribe(message *tgbotapi.Message) error {
	err := b.repo.DeleteSubscriber(message.Chat.ID)
	if err != nil {
		return errUnableToUnsubscribe
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, successfulUnsubscription)
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleCommandCheckSubscribe(message *tgbotapi.Message) error {
	_, err := b.repo.GetSubscriber(message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, subscriptionStatusBad)
		_, err = b.bot.Send(msg)
		if err != nil {
			return err
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, subscriptionStatusGood)
		_, err = b.bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
