package message

import (
	"net/http"
	"schoolonline/routes/access"

	"github.com/gin-gonic/gin"
)

func GetMessageDirectorHandler(c *gin.Context) {

	user := access.GetUser(c)

	c.HTML(http.StatusOK, "message_director.html", user)
}
