package entrance

import (
	"fmt"
	"net/http"

	"schoolonline/dict"
	"schoolonline/webmessage"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationData struct {
	Login    string `form:"login" binding:"required"`
	Password string `form:"password" binding:"required"`
	Error    string `form:"-"`
}

func EntranceGetHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "entrance.html", nil)
}

func EntrancePostHandler(c *gin.Context) {

	var data RegistrationData

	// Привязка данных формы
	if err := c.ShouldBind(&data); err != nil {
		data.Error = err.Error()
		c.HTML(http.StatusOK, "entrance.html", data)
		return
	}

	// Проверка логина пользователя
	user, err := dict.GetUserByLogin(data.Login)
	if err != nil {
		fmt.Println(err)
		data.Error = "неправильный логин, попробуйте еще раз"
		c.HTML(http.StatusOK, "entrance.html", data)
		return
	}

	// Проверка пароля
	if !checkPasswordHash(data.Password, user.Hash) {
		data.Error = "неправильный пароль, попробуйте еще раз"
		c.HTML(http.StatusOK, "entrance.html", data)
		return
	}

	// Сохранение пользователя в сессии при успешной аутентификации
	session := sessions.Default(c)
	session.Set("user", user.Login) // Сохраняем логин пользователя
	if err = session.Save(); err != nil {
		webmessage.Err(c, err, "error saving session:", "/")
		return
	}

	RedirectStartPageByUserRole(c, user.Role)
}

func RedirectStartPageByUserRole(c *gin.Context, role string) {
	switch role {
	case "parent":
		c.Redirect(http.StatusFound, "/start/page/parent")
	case "teacher":
		c.Redirect(http.StatusFound, "/menu")
	case "student":
		c.Redirect(http.StatusFound, "/menu")
	case "director":
		c.Redirect(http.StatusFound, "/menu")
	case "it":
		c.Redirect(http.StatusFound, "/menu")
	default:
		c.Redirect(http.StatusFound, "/menu")
	}
}

// LogoutHandler - обработчик для выхода из аккаунта
func GetLogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/entrance")
}

// LogoutHandler - обработчик для выхода из аккаунта
func PostLogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/entrance")
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
