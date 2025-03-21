package tgbot

import (
	"fmt"
	"schoolonline/tgbot/step"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func parseStep(id int64, ch chan bool, msg *tgbotapi.Message) {

	step, found, err := step.GetStep(id)
	if err != nil || !found {
		ch <- false
		return
	}

	ch <- true

	// step.Print()

	switch step.Function {
	// уведомление о оплате
	case "ar":
		switch step.Step {
		case 1:
			text, err := getMessageText(msg)
			if err != nil {
				fmt.Println(err)
				fmt.Println(err)
			}
			execParentAccauntReplenishmentStep2(id, text, step)
		case 2:
			fmt.Println("ждем фото")
			photo, err := getMessagePhoto(msg)
			if err != nil {
				fmt.Println(err)
				fmt.Println(err)
			}
			fmt.Println(photo.FileID, photo.FileSize)
			execParentAccauntReplenishmentStep3(id, photo.FileID, step)
		}
	}
}
