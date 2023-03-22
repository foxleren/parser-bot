package main

import (
	"github.com/foxleren/parser-bot/pkg/repository"
	"github.com/foxleren/parser-bot/pkg/service/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/siruspen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err := initConfig(); err != nil {
		logrus.Fatalf("Caught error while initializing config: ", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Caught error while loading .env file: ", err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	parserData := telegram.ParserData{
		PythonFile: viper.GetString("parser.pythonfile"),
		ExcelFile:  viper.GetString("parser.excelfile"),
	}

	repos := repository.NewRepository(db)
	//services := service.NewService(repos)

	tgBot := telegram.NewBot(bot, *repos, parserData)
	if err = tgBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
