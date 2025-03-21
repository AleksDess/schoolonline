package access

import (
	"fmt"
	"net/http"
	"schoolonline/dict"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// получаем юзера из сессии
func GetUser(c *gin.Context) dict.User {

	session := sessions.Default(c)
	res := session.Get("user")

	if res == nil {
		c.Redirect(http.StatusFound, "/entrance")
		c.Abort()
		return dict.User{}
	}

	us, err := dict.GetUserByLogin(res.(string))
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/entrance")
		c.Abort()
		return dict.User{}
	}

	return us
}
