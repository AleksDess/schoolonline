package directory

import (
	"net/http"
	"schoolonline/dict"
	"schoolonline/internal"
	"schoolonline/routes/access"
	"schoolonline/webmessage"

	"github.com/gin-gonic/gin"
)

// список школ
func GetDirectoryListSchoolHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	school := dict.GetAllSchoolsFromUser(c, user.Login)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "l_school.html", gin.H{
		"Schools": school,
		"user":    internal.GetUser(user),
	})
}

// удаление школы
func GetDirectoryDeleteSchoolHandler(c *gin.Context) {

	dict.DeleteSchool(c, internal.GetQueryString(c, "id"))
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// добавление школы
func GetDirectoryInputSchoolHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "i_school.html", nil)
}

// добавление школы
func PostDirectoryInputSchoolHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	var sch dict.School
	if err := c.ShouldBind(&sch); err != nil {
		webmessage.Err(c, err, "ошибка чтения формы", "/directory")
		return
	}

	sch.Rec(c, user.Login)
	if c.IsAborted() {
		return
	}

	c.Redirect(http.StatusFound, "/directory")
}

// просмотр школы
func GetDirectoryViewSchoolHandler(c *gin.Context) {

	id := internal.GetQueryString(c, "id")
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "v_school.html", dict.GetSchoolByID(c, id))
}
