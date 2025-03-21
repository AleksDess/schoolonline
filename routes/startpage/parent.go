package startpage

import (
	"net/http"
	"schoolonline/dict"
	"schoolonline/internal"
	"schoolonline/routes/access"

	"github.com/gin-gonic/gin"
)

func GetStartPageHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	res := dict.GetParentInfoByLogin(c, user.Login)
	if c.IsAborted() {
		return
	}

	c.Set("data", res)
	c.Set("user", internal.GetUser(user))
	c.Set("login", user.Login)

	c.HTML(http.StatusOK, "start_parent.html", c.Keys)
}

func PostStartPageHandler(c *gin.Context) {

}
