package directory

import (
	"net/http"
	"schoolonline/internal"
	"schoolonline/kassa"
	"schoolonline/routes/access"

	"github.com/gin-gonic/gin"
)

// list pay
func GetDirectoryListPayHandler(c *gin.Context) {
	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	pay := kassa.Get50KassaViewByRole(c, user.Role, user.Login, user.SchoolID)
	if c.IsAborted() {
		return
	}

	pay.SetOperationName()

	jp := pay.Marshall(c)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "l_pay.html", gin.H{
		"jsonPay": jp,
		"user":    internal.GetUser(user),
	})
}
