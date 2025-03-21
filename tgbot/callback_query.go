package tgbot

import (
	"encoding/json"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallBackPaimentMessage struct {
	D int64 `json:"d"`
	U int64 `json:"u"`
	C int   `json:"c"`
	A bool  `json:"a"`
}

func (a *CallBackPaimentMessage) marshall() string {
	r, _ := json.Marshal(a)
	return string(r)
}

func (a *CallBackPaimentMessage) unMarshall(s string) error {
	err := json.Unmarshal([]byte(s), a)
	return err
}

func sendMessageWithCallbackKeyboardPaimentMessage(dirId, UsId int64, cost int) (err error) {

	msT := CallBackPaimentMessage{D: dirId, U: UsId, C: cost, A: true}
	msF := CallBackPaimentMessage{D: dirId, U: UsId, C: cost, A: false}

	answerT := msT.marshall()
	answerF := msF.marshall()

	// Создание кнопок с callback data
	button1 := tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "cbpm!"+answerT)
	button2 := tgbotapi.NewInlineKeyboardButtonData("Отказать", "cbpm!"+answerF)

	// Создание инлайн-клавиатуры и добавление кнопок
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(button1, button2),
	)

	// Создание и отправка сообщения с клавиатурой
	msg := tgbotapi.NewMessage(dirId, "Выберите опцию:")
	msg.ReplyMarkup = keyboard

	_, err = bot.Send(msg)
	fmt.Println(err)
	return
}

func parseCallbackQuery(callback *tgbotapi.CallbackQuery) {
	// Чтение callback data
	callbackData := callback.Data

	// fmt.Println(callbackData)

	answerList := strings.Split(callbackData, "!")

	if len(answerList) != 2 {
		fmt.Println("error callback data")
		return
	}

	typ := answerList[0]
	answer := answerList[1]

	switch typ {
	case "cbpm":
		go parseCbpmAnswer(answer, callback)
	}

}
