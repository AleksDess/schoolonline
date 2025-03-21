package tgbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getId(msg *tgbotapi.Message) int64 {
	return msg.Chat.ID
}

func getMessageText(msg *tgbotapi.Message) (string, error) {
	if msg == nil || msg.Text == "" {
		return "", fmt.Errorf("текстовое сообщение отсутствует")
	}
	if msg.IsCommand() {
		return "", fmt.Errorf("это команда")
	}
	return msg.Text, nil
}

func getMessageParam(msg *tgbotapi.Message) (string, string, bool, error) {

	if msg == nil || msg.Text == "" {
		return "", "", false, fmt.Errorf("текстовое сообщение отсутствует")
	}

	if !msg.IsCommand() {
		return "", "", false, fmt.Errorf("команда отсутствует")
	}

	command := msg.Command()
	ln := len(msg.CommandArguments())

	if ln == 0 {
		return command, "", false, nil
	}

	param := msg.CommandArguments()

	return command, param, true, nil
}

func getMessagePhoto(msg *tgbotapi.Message) (*tgbotapi.PhotoSize, error) {
	if msg == nil || len(msg.Photo) == 0 {
		return nil, fmt.Errorf("фото отсутствует")
	}
	return &msg.Photo[len(msg.Photo)-1], nil
}

func getMessageVideo(msg *tgbotapi.Message) (*tgbotapi.Video, error) {
	if msg == nil || msg.Video == nil {
		return nil, fmt.Errorf("видео отсутствует")
	}
	return msg.Video, nil
}

func getMessageDocument(msg *tgbotapi.Message) (*tgbotapi.Document, error) {
	if msg == nil || msg.Document == nil {
		return nil, fmt.Errorf("документ отсутствует")
	}
	return msg.Document, nil
}

func GetMessageAudio(msg *tgbotapi.Message) (*tgbotapi.Audio, error) {
	if msg == nil || msg.Audio == nil {
		return nil, fmt.Errorf("аудио отсутствует")
	}
	return msg.Audio, nil
}
