package sendmail

import (
	"fmt"
	"net/smtp"
	"schoolonline/config"
	"schoolonline/webmessage"
	"strings"

	"github.com/gin-gonic/gin"
)

func SendEmailCheckRegistration(c *gin.Context, to, login, password, code string) {

	err := sendEmailCheckRegistration(to, login, password, code)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка отправки письда для верификации почты", "/")
		c.Abort()
		return
	}
}

func sendEmailCheckRegistration(to, login, password, code string) error {

	pass := config.C.EmailKey

	smtpServer := "smtp.gmail.com"
	port := "587"

	from := config.C.EmailSend

	body := createBody(login, password, code)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + from + "\n" + headers + "\n" + body

	auth := smtp.PlainAuth("", from, pass, smtpServer)

	err := smtp.SendMail(smtpServer+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

var template = `
<html>
<body>
    <div class="background-image d-flex align-items-center justify-content-center"></div>
	<p>Вы успешно прошли регистрацию на сайте <strong>smart-crm.org.ua</strong> в качестве владельца онлайн школы</p>

	<p><strong>Логин:  </strong> $1</p>
    <p><strong>Пароль: </strong> $2</p>
    <p><strong>Сайт:   </strong> <a href="http://smart-crm.org.ua">http://smart-crm.org.ua</a></p>

	<p>для подтверждения регистрации пройдите по ссылке ниже</p>
    <p><strong>Подтверждение регистрации: </strong> <a href="$3">$3</a></p>

	<p>на это письмо отвечать не нужно</p>

</body>
</html>
`

func createBody(login, password, code string) string {
	confirmationLink := "https://smart-crm.org.ua/checkemail/" + code
	r := strings.Replace(template, "$1", login, 1)
	r = strings.Replace(r, "$2", password, 1)
	r = strings.ReplaceAll(r, "$3", confirmationLink)
	return r
}

// Access
func SendEmailAccess(c *gin.Context, to, login, password, role string) {

	err := sendEmailAccess(to, login, password, role)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка отправки письда с доступами", "/")
		c.Abort()
		return
	}
}

func sendEmailAccess(to, login, password, role string) error {

	pass := config.C.EmailKey

	smtpServer := "smtp.gmail.com"
	port := "587"

	from := config.C.EmailSend

	body := createBodyAccess(login, password, role)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + from + "\n" + headers + "\n" + body

	auth := smtp.PlainAuth("", from, pass, smtpServer)

	err := smtp.SendMail(smtpServer+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

var template1 = `
<html>
<body>
    <div class="background-image d-flex align-items-center justify-content-center"></div>
	<p>Вы зарегистрированы на сайте <strong>smart-crm.org.ua</strong> в качестве $3</p>

	<p><strong>Логин:  </strong> $1</p>
    <p><strong>Пароль: </strong> $2</p>
    <p><strong>Сайт:   </strong> <a href="http://smart-crm.org.ua">http://smart-crm.org.ua</a></p>
	<p></p>
	<p>на это письмо отвечать не нужно</p>

</body>
</html>
`

func createBodyAccess(login, password, role string) string {
	r := strings.Replace(template1, "$1", login, 1)
	r = strings.Replace(r, "$2", password, 1)
	r = strings.ReplaceAll(r, "$3", role)
	return r
}
