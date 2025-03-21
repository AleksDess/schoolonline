package tgbot

import (
	"encoding/json"
	"fmt"
	"schoolonline/dict"
	"schoolonline/tgbot/step"
	"strconv"
	"time"
)

func parseText(id int64, text string, user dict.User) {

	switch text {

	// director
	case directorText1:
		razrab(id)
	case directorText2:
		razrab(id)

		// parent
	case parentText1:
		execParentAccauntReplenishmentStep1(id, text, user)
	case parentText2:
		getParentBalance(id)
	case parentText3:
		getParentInfo(id)
	case parentText4:
		BotSendTextKeyboard(id, selectAction, generateKeyboard(keyParentRedact, 2))
	case parentText5:
		fmt.Println(id, "отработать: добавить студента")
		razrab(id)
	case parentText6:
		fmt.Println(id, "отработать: удалить студента")
		razrab(id)
	case parentText7:
		fmt.Println(id, "отработать: подписаться на курс")
		razrab(id)
	case parentText8:
		fmt.Println(id, "отработать: отменить курс")
		razrab(id)
	case parentHome:
		BotSendTextKeyboard(id, selectAction, generateKeyboard(keyParentPrimary, 2))

		// teacher
	case teacherText1:
		razrab(id)
	case teacherText2:
		razrab(id)

		// student
	case studentText1:
		razrab(id)
	case studentText2:
		razrab(id)

	default:

	}
}

// инфо родителя
func getParentInfo(id int64) {
	pv, err := dict.GetParentInfoByTgId(id)
	if err != nil {
		fmt.Println(err)
		errMess := `<pre>Не удалось получить информацию счета. 
			Попробуйте позже. 
			Код ошибки: p_info1:1</pre>`
		BotSendTextKeyboard(id, errMess, generateKeyboard(keyParentPrimary, 2))
		return
	}

	mess := fmt.Sprintf(`
	<pre><b>%s</b> 
	<i>на Вашем счете     %d %s</i>
	<i>студентов             %d</i>
	<i>занятий               %d</i>
	<i>занятий в неделю      %d</i>
	<i>цена занятий          %d</i>
	<i>цуна занятий в неделю %d</i>
	</pre>
	`, pv.FullName, int(pv.Balance), pv.CurrencyCod, pv.CountStudent,
		pv.CountLesson, pv.CountLessonWeek, pv.SummaLesson, pv.SummaLessonWeek)
	BotSendTextKeyboard(id, mess, generateKeyboard(keyParentPrimary, 2))

}

// баланс родителя
func getParentBalance(id int64) {
	balance, err := dict.GetUserBalanceByTgId(id)
	if err != nil {
		fmt.Println(err)
		errMess := `<pre>Не удалось получить баланс счета. 
			Попробуйте позже. 
			Код ошибки: p_bal1:1</pre>`
		BotSendTextKeyboard(id, errMess, generateKeyboard(keyParentPrimary, 2))
		return
	}

	mess := fmt.Sprintf(`
	<pre><b>%s</b> 
	<i>на Вашем счете %d %s</i></pre>
	`, balance.Name, balance.Balance, balance.Code)
	BotSendTextKeyboard(id, mess, generateKeyboard(keyParentPrimary, 2))
}

type accauntReplenishment struct {
	IdUser     int64
	Summa      int
	Currency   string
	ParentName string
}

// Функция для преобразования структуры в []byte
func (a *accauntReplenishment) marshall() ([]byte, error) {
	data, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

func accauntReplenishmentUnmarshall(data []byte) (res accauntReplenishment, err error) {

	err = json.Unmarshal(data, &res)
	if err != nil {
		fmt.Println(err)
		return accauntReplenishment{}, err
	}
	return res, nil
}

func execParentAccauntReplenishmentStep1(id int64, text string, user dict.User) (err error) {

	switch text {
	case "пополнение счета":
		fmt.Println("создать степ, ждать сумму")
		info, err := dict.GetParentNameAndCurrencyNameByLogin(user.Login)
		if err != nil {
			fmt.Println(err)
			errMess := `<pre>Не удалось доставить фото оплаты исполнительному директору. 
			Попробуйте передать уведомление об оплате позже.  
			Код ошибки: ar1:1</pre>`
			BotSendParentErrDellStep(id, errMess)
			return err
		}
		ar := accauntReplenishment{IdUser: id, Currency: info.CurrencyCode, ParentName: info.Name}
		arJson, err := ar.marshall()
		if err != nil {
			fmt.Println(err)
			errMess := `<pre>Не удалось передать уведомление об оплате. 
			Попробуйте позже. 
			Код ошибки: ar1:2</pre>`
			BotSendParentErrDellStep(id, errMess)
			return err
		}
		step := step.Step{IdUser: id, Function: "ar", Step: 1, Data: arJson}

		err = step.Rec()
		if err != nil {
			fmt.Println(err)
			errMess := `<pre>Не удалось передать уведомление об оплате. 
			Попробуйте позже. 
			Код ошибки: ar1:3</pre>`
			BotSendParentErrDellStep(id, errMess)
			return err
		}

		mess := fmt.Sprintf("Введите сумму пополнения, %s %s", info.CurrencyCode, info.CurrencyRusName)
		BotSendText(id, mess)
	}
	return nil
}

func execParentAccauntReplenishmentStep2(id int64, text string, step step.Step) {

	ar, err := accauntReplenishmentUnmarshall(step.Data)
	if err != nil {
		fmt.Println(err)
		errMess := `<pre>Не удалось передать уведомление об оплате. 
		Попробуйте позже. 
		Код ошибки: ar2:0</pre>`
		BotSendParentErrDellStep(id, errMess)
		return
	}

	sum, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println(err)
		BotSendText(id, "Введите сумму оплаты (только цифры).")
		return
	}
	ar.Summa = sum

	arJson, err := ar.marshall()
	if err != nil {
		fmt.Println(err)
		errMess := `<pre>Не удалось передать уведомление об оплате. 
		Попробуйте позже.  
		Код ошибки: ar2:1</pre>`
		BotSendParentErrDellStep(id, errMess)
		return
	}
	step.Step = 2
	step.Data = arJson
	err = step.Update()
	if err != nil {
		fmt.Println(err)
		errMess := `<pre>Не удалось передать уведомление об оплате. 
		Попробуйте позже.  
		Код ошибки: ar2:2
			</pre>`
		BotSendParentErrDellStep(id, errMess)
		return
	}

	BotSendText(id, fmt.Sprintf("%s, ожидаем фото или скрин с подтверждением оплаты %d %s", ar.ParentName, ar.Summa, ar.Currency))

}

func execParentAccauntReplenishmentStep3(id int64, photoId string, step step.Step) {

	if step.Function == "ar" {
		dirTgId := dict.GetTgIdDirector(id)
		if dirTgId == 0 {
			errMess := `<pre>Для Вашей школы не найден исполнительный директор. 
			Попробуйте передать уведомление об оплате позже.  
			Код ошибки: ar3:0</pre>`
			BotSendParentErrDellStep(id, errMess)
			return
		}

		ar, err := accauntReplenishmentUnmarshall(step.Data)
		if err != nil {
			fmt.Println(err)
			errMess := `<pre>Не удалось передать уведомление об оплате. 
			Попробуйте позже. 
			Код ошибки: ar3:1
			</pre>`
			BotSendParentErrDellStep(id, errMess)
			return
		}

		mess := fmt.Sprintf("фото оплаты от %s на сумму %d %s", ar.ParentName, ar.Summa, ar.Currency)

		err = BotSendPhoto(dirTgId, photoId, mess)
		if err != nil {
			fmt.Println(err)
			errMess := `<pre>Не удалось доставить фото оплаты исполнительному директору. 
			Попробуйте передать уведомление об оплате позже.  
			Код ошибки: ar3:2
			</pre>`
			BotSendParentErrDellStep(id, errMess)
			return
		}

		err = sendMessageWithCallbackKeyboardPaimentMessage(dirTgId, id, ar.Summa)
		if err != nil {
			fmt.Println(err)
			errMess := `<pre>Не удалось доставить исполнительную клавиатуру  исполнительному директору. 
			Попробуйте передать уведомление об оплате позже.  
			Код ошибки: ar3:3
			</pre>`
			BotSendParentErrDellStep(id, errMess)
			return
		}

		mess = fmt.Sprintf(`<pre>"%s, 
		фото или скрин с подтверждением оплаты %d %s 
		доставлено успешно."
			</pre>`, ar.ParentName, ar.Summa, ar.Currency)
		BotSendText(id, mess)
		<-time.After(500 * time.Millisecond)

		mess = fmt.Sprintf(`<pre>"%s, 
		ожидайте подтверждения оплаты</pre>`, ar.ParentName)
		BotSendTextKeyboard(id, mess, generateKeyboard(keyParentPrimary, 2))

	}
	// передача калбек клавиатуры

	step.Delete(id)

}
