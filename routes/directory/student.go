package directory

import (
	"fmt"
	"net/http"
	"schoolonline/crypto"
	"schoolonline/dict"
	"schoolonline/internal"
	"schoolonline/routes/access"
	"time"

	"github.com/gin-gonic/gin"
)

// список student
func GetDirectoryListStudentHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	var stud = dict.GetListStudentByUserRole(c, user)
	var lpar = dict.ListParent{}

	switch user.Role {
	case "it", "director":
		lpar = dict.GetAllParentFromUser(c, user.Login)
		if c.IsAborted() {
			return
		}
	case "parent":
		lpar = dict.GetParentFromUser(c, user.Login)
		if c.IsAborted() {
			return
		}
	}

	c.Set("students", stud.Marshall(c))
	c.Set("parents", lpar.Marshall(c))
	c.Set("user", internal.GetUser(user))

	c.HTML(http.StatusOK, "l_student.html", c.Keys)
}

// удаление student
func GetDirectoryDeleteStudentHandler(c *gin.Context) {

	id := internal.GetQueryString(c, "id")
	if c.IsAborted() {
		return
	}

	dict.DeleteStudent(c, id)
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ---------------------------------------------------
// добавление student
func PostDirectoryInputStudentHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	st := dict.User{}

	st.FirstName = internal.GetFormaString(c, "first_name")
	if c.IsAborted() {
		return
	}
	st.LastName = internal.GetFormaString(c, "last_name")
	if c.IsAborted() {
		return
	}
	st.YearOfBirth = internal.GetFormaInt(c, "year_of_birth")
	if c.IsAborted() {
		return
	}
	st.ParentLogin = internal.GetFormaString(c, "parent_login")
	if c.IsAborted() {
		return
	}

	currency, _ := dict.GetCurrencyByUserLogin(st.ParentLogin)

	st.Role = "student"
	st.PassWord = crypto.GeneratePassword()
	st.Login = crypto.GetUlid()
	st.Creator = user.Login
	st.SchoolID = user.SchoolID
	st.CurrencyId = currency.ID
	st.CreateTime = time.Now()
	st.VerifyEmail = true
	st.UpdateTime = time.Now()
	st.VerifyTime = time.Now()

	fmt.Println(st)

	st.Rec(c)
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Студент добавлен",
	})
}

// просмотр student
func GetDirectoryViewStudentHandler(c *gin.Context) {

	login := internal.GetQueryString(c, "login")
	if c.IsAborted() {
		return
	}

	student := dict.GetStudentViewById(c, login)
	if c.IsAborted() {
		return
	}

	student.ListLesson = dict.GetAllLessonViewByStudentLogin(c, login)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "v_student.html", student)
}
