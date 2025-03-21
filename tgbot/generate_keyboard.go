package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func generateKeyboard(line []string, col int) tgbotapi.ReplyKeyboardMarkup {
	// Создаем клавиатуру
	var keyboard [][]tgbotapi.KeyboardButton

	ln := len(line)

	if col == 1 {
		for _, i := range line {
			button := tgbotapi.NewKeyboardButton(i)
			row := []tgbotapi.KeyboardButton{}
			row = append(row, button)
			keyboard = append(keyboard, row)
		}

	} else if col == 2 {
		// Перебираем элементы строки и создаем кнопки
		for i := 0; i < ln; {
			row := []tgbotapi.KeyboardButton{}

			// Добавляем по 2 кнопки, если это возможно
			if i+1 < ln {
				button1 := tgbotapi.NewKeyboardButton(line[i])
				button2 := tgbotapi.NewKeyboardButton(line[i+1])
				row = append(row, button1, button2)
				i += 2
			} else {
				// Если не хватает второй кнопки, добавляем одну
				button := tgbotapi.NewKeyboardButton(line[i])
				row = append(row, button)
				i++
			}

			// Добавляем сформированную строку кнопок в клавиатуру
			keyboard = append(keyboard, row)
		}

	}

	// Возвращаем клавиатуру
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}
}
