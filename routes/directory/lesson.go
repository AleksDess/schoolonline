package directory

import (
	"fmt"
	"net/http"
	"schoolonline/dict"
	"schoolonline/internal"
	"schoolonline/routes/access"
	"schoolonline/transaction"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

// список уроков
func GetListLessonHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	var lessons dict.ListLessonView

	switch user.Role {
	case "director":
		lessons = dict.GetAllLessonViewByUserLoginDirector(c, user.Login)
		if c.IsAborted() {
			return
		}
	case "parent":
		lessons = dict.GetAllLessonViewByUserLoginParent(c, user.Login)
		if c.IsAborted() {
			return
		}
	case "teacher":
		lessons = dict.GetAllLessonViewByUserLoginTeacher(c, user.Login)
		if c.IsAborted() {
			return
		}
	case "it":
		lessons = dict.GetAllLessonViewByUserLoginIt(c)
		if c.IsAborted() {
			return
		}
	}

	c.Set("lessons", lessons.Marshall(c))
	c.Set("user", internal.GetUserJson(c, user))
	c.HTML(http.StatusOK, "l_lessons.html", c.Keys)
}

// просмотр урока
func GetViewLessonHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	lesson := dict.GetLessonViewById(c)
	if c.IsAborted() {
		return
	}

	us := fmt.Sprintf("{\"user_role\": \"%s\"}", user.Role)

	c.Set("lesson", lesson.Marshall(c))
	c.Set("user", us)

	c.HTML(http.StatusOK, "v_lesson.html", c.Keys)
}

// добавление урока гет
func GetDirectoryInputLessonHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	type Data struct {
		ListItem    dict.ListItem
		ListTeacher dict.ListTeacher
		ListStudent dict.ListStudent
		SchoolID    int
		UserRole    string
		UrlRedirect string
	}

	var stud = dict.ListStudent{}

	switch user.Role {
	case "it", "director":
		stud = dict.GetListStudentByUser(c, user.Login)
		if c.IsAborted() {
			return
		}
	case "parent":
		stud = dict.GetListStudentByParent(c, user.Login)
		if c.IsAborted() {
			return
		}
	case "teacher":
		stud = dict.GetListStudentByTeacher(c, user.Login)
		if c.IsAborted() {
			return
		}
	default:
		webmessage.SendMessage(c, "У вас нет доступа к списку студентов", "/menu")
		return
	}

	rs := Data{}
	rs.ListItem = dict.GetAllItemsBySchoolId(c, fmt.Sprint(user.SchoolID))
	rs.ListTeacher = dict.GetAllTeachersBySchoolId(c, fmt.Sprint(user.SchoolID))
	rs.ListStudent = stud
	rs.SchoolID = user.SchoolID
	rs.UserRole = user.Role
	rs.UrlRedirect = "get"

	c.HTML(http.StatusOK, "i_lessons.html", rs)

}

// добавление урока пост
func PostDirectoryInputLessonHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	studentLogin := internal.GetFormaString(c, "student_login")
	if c.IsAborted() {
		return
	}

	type Data struct {
		ListItem           dict.ListItem
		ListTeacher        dict.ListTeacher
		ListStudent        dict.ListStudent
		SchoolID           string
		StudentLogin       string
		ParentLogin        string
		CurrencyParentID   string
		CurrencyParentName string
		UserRole           string
		UrlRedirect        string
	}

	rs := Data{}
	rs.ListItem = dict.GetAllItemsBySchoolId(c, fmt.Sprint(user.SchoolID))
	if c.IsAborted() {
		return
	}
	rs.ListTeacher = dict.GetAllTeachersBySchoolId(c, fmt.Sprint(user.SchoolID))
	if c.IsAborted() {
		return
	}

	rs.ListStudent = dict.GetListStudentByStudentLogin(c, studentLogin)
	if c.IsAborted() {
		return
	}
	rs.SchoolID = fmt.Sprint(user.SchoolID)
	rs.UserRole = user.Role

	rs.UrlRedirect = "post"

	c.HTML(http.StatusOK, "i_lessons.html", rs)

}

// запись урока
func PostDirectorySaveLessonHandler(c *gin.Context) {

	rs := dict.Lesson{}

	if err := c.ShouldBind(&rs); err != nil {
		fmt.Println(err)
		return
	}

	typLesson := internal.GetFormaString(c, "lesson_type")
	if c.IsAborted() {
		return
	}

	switch typLesson {
	case "individual":
		rs.IsGroup = false
	case "group":
		rs.IsGroup = true
	}

	rs.CreateTime = time.Now()
	rs.CloseTime = time.Date(2099, 12, 31, 23, 59, 0, 0, time.UTC)

	// rs.Print()

	rs.Rec(c)
	if c.IsAborted() {
		return
	}

	redirect := c.PostForm("url_redirect")

	switch redirect {
	case "post":
		strReturn := fmt.Sprintf("/view/student?login=%s", rs.StudentLogin)
		c.Redirect(http.StatusFound, strReturn)
	case "get":
		c.Redirect(http.StatusFound, "/list/lesson")
	default:
		c.Redirect(http.StatusFound, "/menu")
	}

}

// отмена урока
func GetLessonCancelHandler(c *gin.Context) {

	id := dict.UpdateLessonClose(c)
	if c.IsAborted() {
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/view/lesson?id=%s", id))
}

// проведение урока
// GetLessonCompleteHandler - обработка проведения урока (AJAX)
func GetLessonCompleteHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		c.JSON(http.StatusOK, gin.H{"success": false, "cause": 1})
		return
	}

	lesson := dict.GetLessonByIdWeb(c)
	if c.IsAborted() {
		c.JSON(http.StatusOK, gin.H{"success": false, "cause": 2})
		return
	}

	// Проверка, отмечен ли урок уже сегодня
	if lesson.CheckLessonCompletionDate() {
		c.JSON(http.StatusOK, gin.H{"success": false, "cause": 3})
		return
	}

	// Выполнение транзакции для проведения урока
	transaction.ExecTransactionLessonComplete(c, lesson, user)
	if c.IsAborted() {
		c.JSON(http.StatusOK, gin.H{"success": false, "cause": 4})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"success": true, "cause": 0})
}
