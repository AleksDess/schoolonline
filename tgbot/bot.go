package tgbot

import (
	"fmt"
	"log"
	"schoolonline/config"
	"schoolonline/launch"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var err error

const workerCount = 5
const lenUpdareChan = 100

func RunBot() {

	botkey := config.C.BotKey

	if launch.Launch == "home" {
		botkey = config.C.BotKeyTest
	}

	bot, err = tgbotapi.NewBotAPI(botkey)
	if err != nil {
		fmt.Println(err)
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	updateChan := make(chan tgbotapi.Update, lenUpdareChan)

	for i := 0; i < workerCount; i++ {
		go worker(updateChan)
	}

	for update := range updates {

		select {
		case updateChan <- update:
		default:
			log.Println("Update channel buffer is full")
		}
	}
}
