package tgbot

import (
	"schoolonline/tgbot/step"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BotSendText(id int64, s string) error {
	mess := tgbotapi.NewMessage(id, s)
	mess.ParseMode = "HTML"
	_, err := bot.Send(mess)
	return err
}

func BotSendPhoto(id int64, fileID string, caption string) error {
	photo := tgbotapi.NewPhoto(id, tgbotapi.FileID(fileID))
	photo.Caption = caption
	photo.ParseMode = "HTML"
	_, err := bot.Send(photo)
	return err
}

func BotSendVideo(id int64, videoPath string, caption string) error {
	video := tgbotapi.NewVideo(id, tgbotapi.FilePath(videoPath))
	video.Caption = caption
	video.ParseMode = "HTML"
	_, err := bot.Send(video)
	return err
}

func BotSendAudio(id int64, audioPath string, caption string) error {
	audio := tgbotapi.NewAudio(id, tgbotapi.FilePath(audioPath))
	audio.Caption = caption
	audio.ParseMode = "HTML"
	_, err := bot.Send(audio)
	return err
}

func BotSendDocument(id int64, docPath string, caption string) error {
	document := tgbotapi.NewDocument(id, tgbotapi.FilePath(docPath))
	document.Caption = caption
	document.ParseMode = "HTML"
	_, err := bot.Send(document)
	return err
}

func BotSendTextKeyboard(id int64, s string, keyboard tgbotapi.ReplyKeyboardMarkup) error {
	// Создаем сообщение
	mess := tgbotapi.NewMessage(id, s)
	mess.ParseMode = "HTML"

	// Добавляем клавиатуру к сообщению
	mess.ReplyMarkup = keyboard

	// Отправляем сообщение с клавиатурой
	_, err := bot.Send(mess)
	return err
}

func BotSendParentErrDellStep(id int64, text string) {
	BotSendText(id, text)
	<-time.After(500 * time.Millisecond)
	BotSendTextKeyboard(id, selectAction, generateKeyboard(keyParentPrimary, 2))
	step.Delete(id)
}
