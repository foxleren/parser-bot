package telegram

import (
	"errors"
	"fmt"
	"github.com/foxleren/parser-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/siruspen/logrus"
	"log"
	"os"
	"os/exec"
	"time"
)

type Bot struct {
	bot        *tgbotapi.BotAPI
	repo       repository.Repository
	parserData ParserData
}

type ParserData struct {
	PythonFile string
	ExcelFile  string
}

func NewBot(bot *tgbotapi.BotAPI, repo repository.Repository, data ParserData) *Bot {
	return &Bot{
		bot:        bot,
		repo:       repo,
		parserData: data,
	}
}

var parsingTime string

func setParsingTime() {
	t := time.Now()
	parsingTime = fmt.Sprintf("%d-%02d-%02d %02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute())
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	setParsingTime()

	updates := b.initUpdatesChannel()
	go b.sendMessageToSubscribers()
	err := b.handleUpdates(updates)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			err := b.handleCommand(update.Message)
			if err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) sendMessageToSubscribers() {
	for {
		b.compileParser()

		subscribers, err := b.repo.GetAllSubscribers()
		log.Printf("%v", err)
		for _, sbs := range subscribers {
			go b.sendMessage(sbs.ChatId)
		}

		time.Sleep(24 * time.Hour)
	}
}

func (b *Bot) compileParser() {
	cmd := exec.Command("python3", b.parserData.PythonFile)
	err := cmd.Run()
	if err != nil {
		logrus.Printf("Error while running python code")
	}

	setParsingTime()
}

func (b *Bot) sendMessage(chatId int64) {
	filePath := b.parserData.ExcelFile

	file, err := os.Open(filePath)
	if !errors.Is(err, os.ErrNotExist) {
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Panic(err)
		}

		fileName := fileInfo.Name()

		fileBytes := make([]byte, fileInfo.Size())

		_, err = file.Read(fileBytes)
		if err != nil {
			log.Panic(err)
		}

		fileBytesConfig := tgbotapi.FileBytes{Name: fileName, Bytes: fileBytes}

		msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("Версия от %s", parsingTime))
		b.bot.Send(msg)

		doc := tgbotapi.NewDocument(chatId, fileBytesConfig)
		b.bot.Send(doc)
	}
}
