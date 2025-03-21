package tgbot

import (
	"fmt"
	"schoolonline/dict"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// worker обрабатывает обновления из канала
func worker(updateChan <-chan tgbotapi.Update) {
	for update := range updateChan {

		// если каллбек кнопка
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			go parseCallbackQuery(callback)
			continue
		}

		// если сообщение
		if update.Message != nil {

			msg := update.Message

			// Создаем канал для обмена решениями
			ch := make(chan bool)
			id := getId(msg)

			isUser := true

			// проверяем существование юзера
			user, err := dict.GetUserByTgId(id)
			if err != nil {
				fmt.Println(err)
				isUser = false
			}

			// проверяем сообщение с командой
			command, param, isParam, err := getMessageParam(msg)
			if err == nil {
				// обработка старта с параметром
				if isParam {
					go parseCommand(id, command, param)
					continue
				} else {
					if isUser {
						switch user.Role {
						case "director":
							BotSendTextKeyboard(id, selectAction, generateKeyboard(keyDirectorPrimary, 2))
						case "parent":
							BotSendTextKeyboard(id, selectAction, generateKeyboard(keyParentPrimary, 2))
						case "teacher":
							BotSendTextKeyboard(id, selectAction, generateKeyboard(keyTeacherPrimary, 2))
						case "student":
							BotSendTextKeyboard(id, selectAction, generateKeyboard(keyStudentPrimary, 2))
						}
					} else {
						BotSendText(id, "Вы не зарегистрированы пользователем smart-crm.org.ua.")
					}
					continue
				}
			}

			// проверяем если юзер на шаге
			go parseStep(id, ch, msg)
			fmt.Println("bot message id:", id)

			if <-ch {
				continue
			}

			text, err := getMessageText(msg)
			if err == nil {
				// work photo
				fmt.Println("--text-- :", text)
				go parseText(id, text, user)
			}

			photo, err := getMessagePhoto(msg)
			if err == nil {
				// work photo
				fmt.Println("**photo** id:", photo.FileID)
			}

			video, err := getMessageVideo(msg)
			if err == nil {
				// work video
				fmt.Println("**video** id:", video.FileID)
			}

			audio, err := GetMessageAudio(msg)
			if err == nil {
				// work audio
				fmt.Println("**audio** id:", audio.FileID)
			}

			document, err := getMessageDocument(msg)
			if err == nil {
				// work document
				fmt.Println("**document** id:", document.FileID)
			}
		} else {
			continue
		}
	}
}
