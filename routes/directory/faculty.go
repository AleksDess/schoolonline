package directory

import (
	"net/http"
	"schoolonline/dict"
	"schoolonline/internal"
	"schoolonline/routes/access"
	"schoolonline/webmessage"

	"github.com/gin-gonic/gin"
)

// list faculty
func GetDirectoryListFacultyHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	facultys := dict.GetAllFacultysFromUser(c)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "l_faculty.html", gin.H{
		"Facultys": facultys,
		"user":     internal.GetUser(user),
	})
}

// удаление факультета
func GetDirectoryDeleteFacultyHandler(c *gin.Context) {

	dict.DeleteFaculty(c)
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// добавление факультета
func GetDirectoryInputFacultyHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "i_faculty.html", nil)
}

// -------------------------------------------------------------
// добавление факультета
func PostDirectoryInputFacultyHandler(c *gin.Context) {

	user := access.GetUser(c)

	data := dict.Faculty{}

	if err := c.ShouldBind(&data); err != nil {
		webmessage.Err(c, err, "ошибка чтения формы факультета", "/directory")
		return
	}

	sm := dict.GetIdSchoolByUser(c, user.Login)
	if c.IsAborted() {
		return
	}

	data.School = sm.ID

	data.Rec(c)
	if c.IsAborted() {
		return
	}

	c.Redirect(http.StatusFound, "/list/faculty")
}

// просмотр факультета
func GetDirectoryViewFacultyHandler(c *gin.Context) {

	faculty := dict.GetFacultyByID(c)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "v_faculty.html", faculty)
}
