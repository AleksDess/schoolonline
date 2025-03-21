package registration

import (
	"fmt"
	"schoolonline/dict"
	"schoolonline/webmessage"

	"github.com/gin-gonic/gin"
)

func ConfirmationOfRegisterPostHandler(c *gin.Context) {

	code := c.Param("param")

	user, err := dict.GetUserByVerifyCode(code)

	if err != nil {
		fmt.Println(err)
		webmessage.SendMessage(c, "данная ссылка больше не активна", "/menu")
		return
	}

	if user.VerifyEmail {
		webmessage.SendMessage(c, "Вы ранее уже успешно осуществили подтверждение почтового адреса", "/menu")
		return
	}

	err, verif := dict.UpdateUserVerifyStatus(code)

	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка записи верификации почты в БД", "/menu")
		return
	}

	if verif {
		webmessage.SendMessage(c, "Вы успешно осуществили подтверждение почтового адреса", "/menu")
	} else {
		webmessage.Err(c, err, "ошибка записи верификации почты в БД", "/menu")
	}
}
