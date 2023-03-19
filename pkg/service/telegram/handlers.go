package telegram

import (
	"github.com/foxleren/parser-bot/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/siruspen/logrus"
	"log"
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
	failedUnsubscription     = "Не удалось отписаться от рассылки."
	subscriptionStatusGood   = "Статус подписки: активирована."
	subscriptionStatusBad    = "Статус подписки: деактивирована."

	subscribeError = "Не удалось подписаться на рассылку либо вы уже подписаны. Для проверки подписки используйте команду \n/check_subscribe."
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleCommandStart(message)
	case commandSubscribe:
		return b.handleCommandSubscribe(message)
	case commandCheckSubscribe:
		return b.handleCommandCheckSubscribe(message)
	case commandUnsubscribe:
		return b.handleCommandUnsubscribe(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID

	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleCommandStart(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, replyStart)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleCommandSubscribe(message *tgbotapi.Message) error {
	subscriber := models.Subscriber{ChatId: message.Chat.ID}

	var id int
	id, bd_err := b.repo.CreateSubscriber(subscriber)
	if bd_err != nil {
		logrus.Printf("%v", bd_err)
		msg := tgbotapi.NewMessage(message.Chat.ID, subscribeError)
		_, err := b.bot.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, successfulSubscription)
	_, err := b.bot.Send(msg)

	logrus.Printf("Created id: %d", id)

	b.sendMessage(message.Chat.ID)

	return err
}

func (b *Bot) handleCommandCheckSubscribe(message *tgbotapi.Message) error {
	_, err := b.repo.GetSubscriber(message.Chat.ID)
	if err != nil {
		logrus.Printf("%v", err)

		msg := tgbotapi.NewMessage(message.Chat.ID, subscriptionStatusBad)
		b.bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, subscriptionStatusGood)
		b.bot.Send(msg)
	}

	return nil
}

func (b *Bot) handleCommandUnsubscribe(message *tgbotapi.Message) error {
	err := b.repo.DeleteSubscriber(message.Chat.ID)
	if err != nil {
		logrus.Printf("%v", err)

		msg := tgbotapi.NewMessage(message.Chat.ID, failedUnsubscription)
		b.bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, successfulUnsubscription)
		b.bot.Send(msg)
	}

	return nil
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, replyUnknownCommand)
	_, err := b.bot.Send(msg)

	return err
}
