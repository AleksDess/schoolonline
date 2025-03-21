package tgbot

import (
	"fmt"
	"log"
	"schoolonline/transaction"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func parseCbpmAnswer(answer string, callback *tgbotapi.CallbackQuery) {

	cbpm := CallBackPaimentMessage{}
	err := cbpm.unMarshall(answer)
	if err != nil {
		fmt.Println(err)
		return
	}

	responseText := ""
	if cbpm.A {

		responseText = "Вы подтвердили оплату"

		// сделать зачисление суммы на баланс юзера
		// сделать запись в kassa
		err := transaction.TransactionUserKassa(cbpm.D, cbpm.U, cbpm.C)
		if err != nil {
			fmt.Println(err)
			BotSendText(cbpm.U, "Ошибка записи платежа в БД."+fmt.Sprint(err))
			return
		}

		// сделать сообщение юзеру о зачислении денег
		BotSendText(cbpm.U, "Ваш платеж успешно обработан и зачислен.")
	} else {
		// fmt.Println("no complete")
		responseText = "Вы отклонили оплату"

		// сделать запись юзеру об отклонении оплаты
		BotSendText(cbpm.U, "Ваш платеж отклонен.")
	}

	workingCallBackQueryMessage(callback, responseText)
}

func workingCallBackQueryMessage(callback *tgbotapi.CallbackQuery, responseText string) {
	// Удаление клавиатуры
	editMarkup := tgbotapi.NewEditMessageReplyMarkup(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
	)

	if _, err := bot.Request(editMarkup); err != nil {
		log.Println("Ошибка удаления клавиатуры:", err)
	}

	// Ответ на callback-запрос для удаления индикатора загрузки
	answerNew := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(answerNew); err != nil {
		log.Println("Ошибка обработки callback-запроса:", err)
	}

	// Ответ пользователю на выбор кнопки
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, responseText)
	if _, err := bot.Send(msg); err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}
