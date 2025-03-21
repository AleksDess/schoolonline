package registration

import (
	"fmt"
	"net/http"
	"schoolonline/sendmail"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

func CapchaGetHandler(c *gin.Context) {
	cap := generateCapcha()
	c.HTML(http.StatusOK, "capcha.html", cap)
}

func CapchaPostHandler(c *gin.Context) {

	res := c.PostForm("isHuman")

	if res != "true" {
		webmessage.Err(c, nil, "ошибка проверки", "/")
		return
	}

	session := sessions.Default(c)
	data, ok := session.Get("registrationData").(RegistrationData)
	if !ok {
		webmessage.Err(c, nil, "ошибка получения данных из сессии", "/")
		return
	}

	r := data.createUser()

	// Генерация UUID для ссылки подтверждения
	r.VerifyCode = uuid.New().String()

	r.Rec(c)
	if c.IsAborted() {
		return
	}

	sendmail.SendEmailCheckRegistration(c, r.Email, r.Login, r.PassWord, r.VerifyCode)
	if c.IsAborted() {
		return
	}

	// Сохраняем пользователя в сессии для автоматического входа
	session.Set("user", r.Login)
	session.Save()

	ms := webmessage.MessageList{}
	ms.Titul = fmt.Sprintf("Поздравляем %s", r.Login)
	ms.Message = []string{
		"Регистрация успешно произведена.",
		fmt.Sprintf("На Вашу почту %s отправлено письмо для подтверждения.", r.Email),
		"Сейчас Вам нужно добавить Вашу школу.",
	}
	ms.Redirect = "/input/school"
	ms.Close = false

	ms.SendMessageList(c)
}

type CapchaData struct {
	Num1        int `form:"num1" binding:"required"`
	Num2        int `form:"num2" binding:"required"`
	Num3        int `form:"num3" binding:"required"`
	Answer      int `form:"answer" binding:"required"`
	GivenAnswer int `form:"givenanswer" binding:"required"`
}

func generateCapcha() (r CapchaData) {
	rand.Seed(uint64(time.Now().UnixNano())) // Инициализация генератора случайных чисел

	r.Num1 = rand.Intn(90) + 10 // Генерация числа от 10 до 99
	r.Num2 = rand.Intn(90) + 10 // Генерация числа от 10 до 99
	r.Num3 = rand.Intn(90) + 10 // Генерация числа от 10 до 99

	r.GivenAnswer = r.Num1 + r.Num2 + r.Num3
	r.Answer = 0
	return
}
