package tgbot

func razrab(id int64) {
	BotSendText(id, "Функция в стадии разработки")
}

const selectAction = "выберите действие"
const parentHome = "на главную"

// const parentBack = "назад"
// const parentCancel = "отменить"
const parentText1 = "пополнение счета"
const parentText2 = "баланс счета"
const parentText3 = "информация об аккаунте"
const parentText4 = "редактировать аккаунт"
const parentText5 = "добавить студента"
const parentText6 = "удалить студента"
const parentText7 = "подписаться на курс"
const parentText8 = "отменить курс"

// const parentText9 =
// const parentText10 =
// const parentText11 =

const directorText1 = "статистика"
const directorText2 = "непроверенные поступления"

const teacherText1 = "индивидуальное занятие завершено"
const teacherText2 = "групповое занятие завершено"

const studentText1 = "мои уроки"
const studentText2 = "отменить урок"
