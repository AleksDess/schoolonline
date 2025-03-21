package registration

import (
	"database/sql"
	"fmt"
	"net/http"
	"schoolonline/dict"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type RegistrationData struct {
	Login         string `form:"login"`
	Email         string `form:"email"`
	Phone         string `form:"phone"`
	Password      string `form:"password"`
	Role          string `form:"role"`
	FirstName     string `form:"first_name"`
	LastName      string `form:"last_name"`
	PimentDetails string `form:"piment_details"`
	Error         string `form:"-"`
}

func (a *RegistrationData) createUser() dict.User {
	r := dict.User{}
	r.Login = a.Login
	r.PassWord = a.Password
	r.Role = a.Role
	r.Email = a.Email
	r.Phone = a.Phone
	r.Creator = a.Login
	r.PimentDetails = a.PimentDetails
	r.FirstName = a.FirstName
	r.LastName = a.LastName
	r.CreateTime = time.Now()
	r.UpdateTime = time.Now()
	return r
}

func RegisterGetHandler(c *gin.Context) {

	type data struct {
		Role string
	}

	d := data{Role: "director"}

	c.HTML(http.StatusOK, "registration.html", d)
}

func RegisterPostHandler(c *gin.Context) {

	var data RegistrationData
	if err := c.ShouldBind(&data); err != nil {
		webmessage.Err(c, err, "ошибка контекста регитрации", "/")
		return
	}

	session := sessions.Default(c)
	session.Set("registrationData", data)
	err := session.Save()
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка сессии регистрации", "/")
		return
	}

	c.Redirect(http.StatusFound, "/capcha")
}

func CheckUsernameHandler(c *gin.Context) {

	username := c.Query("login")
	exists := false
	_, err := dict.GetUserByLogin(username)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows { // если пользователь не найден
			exists = false
		} else {
			// обработка других ошибок, например ошибка подключения к базе данных
			c.JSON(500, gin.H{
				"error": "Ошибка сервера",
			})
			return
		}
	} else {
		exists = true
	}
	c.JSON(200, gin.H{
		"exists": exists,
	})
}
