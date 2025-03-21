package tgbot

import (
	"fmt"
	"schoolonline/dict"
)

func parseCommand(id int64, s, param string) {

	if s != "" {
		// регистрация Id в телеграмм боте
		if s == "start" {
			if param != "" {
				fmt.Println("parameter start bot:", param)
				if err != nil {
					fmt.Println(err)
					err = BotSendText(id, "Не распознан параметр команды")
					if err != nil {
						fmt.Println(err)
						fmt.Println("tgBot err:", err)
						return
					}
					return
				}

				// проверяем существование юзера
				user, _ := dict.GetUserByTgId(id)

				if user.TgId == id {
					err = BotSendText(id, "Вы уже зарегистрированы в smart-crm бот.")
					if err != nil {
						fmt.Println(err)
						fmt.Println("tgBot err:", err)
						return
					}
				} else {

					err = dict.UpdateUserTgId(param, id)
					if err != nil {
						fmt.Println(err)
						err = BotSendText(id, "ошибка обновления TgId пользователя")
						if err != nil {
							fmt.Println(err)
							fmt.Println("tgBot err:", err)
							return
						}
					}

					var ms string
					if user.TgId == 0 {
						ms = fmt.Sprintf(`
						Добрый день %s
						вы зарегистрированы в smart-crm бот.`, param)
					} else {
						ms = fmt.Sprintf(`
						Добрый день %s
						ваши данные обновлены в smart-crm бот.`, param)
					}

					err = BotSendText(id, ms)
					if err != nil {
						fmt.Println(err)
						fmt.Println("tgBot err:", err)
						return
					}
				}
			}
		}
	}
}
