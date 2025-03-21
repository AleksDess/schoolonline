package directory

import (
	"fmt"
	"net/http"
	"schoolonline/crypto"
	"schoolonline/dict"
	"schoolonline/internal"
	"schoolonline/routes/access"
	"schoolonline/tgbot"
	"schoolonline/transaction"
	"schoolonline/webmessage"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// список учителей
func GetDirectoryListTeacherHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	teachers := dict.GetAllTeacherFromUser(c, user.Login)
	if c.IsAborted() {
		return
	}

	c.Set("teachers", teachers.Marshall(c))
	c.Set("schoolId", user.SchoolID)
	c.Set("user", internal.GetUserJson(c, user))

	c.HTML(http.StatusOK, "l_teacher.html", c.Keys)
}

// удаление учителя
func GetDirectoryDeleteTeacherHandler(c *gin.Context) {

	id := internal.GetQueryString(c, "id")
	if c.IsAborted() {
		return
	}

	dict.DeleteTeacher(c, id)
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ---------------------------------------------------
// добавление учителя пост запрос
func PostDirectoryInputTeacherHandler(c *gin.Context) {

	isUser := true

	user := access.GetUser(c)
	if c.IsAborted() {
		isUser = false
	}

	teacher := dict.User{}

	if err := c.ShouldBind(&teacher); err != nil {
		webmessage.Err(c, err, "Ошибка чтения формы", "/directory")
		c.Abort()
		return
	}

	us := dict.User{}
	us.Login = c.PostForm("login")
	us.PassWord = c.PostForm("password")
	us.Role = "teacher"
	if isUser {
		us.Creator = user.Login
	} else {
		us.Creator = us.Login
	}
	us.Email = c.PostForm("email")
	us.FirstName = c.PostForm("first_name")
	us.LastName = c.PostForm("last_name")
	us.SchoolID, _ = strconv.Atoi(c.PostForm("school_id"))
	us.Phone = c.PostForm("phone")
	us.CurrencyId, _ = strconv.Atoi(c.PostForm("currency_id"))
	us.PimentDetails = c.PostForm("piment_details")
	us.CreateTime = time.Now()
	us.UpdateTime = time.Now()
	us.VerifyEmail = true
	us.VerifyTime = time.Now()

	typeRegistration := c.PostForm("type_registration")
	linkCode := c.PostForm("link_code")

	switch typeRegistration {
	case "manual":
		transaction.ExecTransactionInsertUserAndTeacher(c, us)
		if c.IsAborted() {
			return
		}
		c.Redirect(http.StatusFound, "/list/teacher")
	case "link":
		transaction.ExecTransactionInsertUserAndTeacherAndUpdateLink(c, us, linkCode)
		if c.IsAborted() {
			return
		}
		c.Redirect(http.StatusFound, "/menu")
	}
}

// добавление учителя вручную гет запрос
func GetDirectoryInputTeacherHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	var data struct {
		SchoolId         int
		SchoolName       string
		TypeRegistration string
		LinkCode         string
		ListCurrency     dict.ListCurrency
		Password         string
	}

	data.SchoolId = user.SchoolID
	data.TypeRegistration = "manual"
	data.LinkCode = "_"

	data.ListCurrency = dict.GetAllCurrencyBySchoolId(c, data.SchoolId)

	data.Password = crypto.GeneratePassword()

	c.HTML(http.StatusOK, "i_teacher.html", data)
}

// добавление учителя по ссылке гет запрос
func GetDirectoryInputTeacherRegistration(c *gin.Context) {

	param := c.Param("param")

	slise := strings.Split(param, ":")

	if len(slise) != 2 {
		webmessage.Err(c, fmt.Errorf("invalid paraveters"), "неправильный параметра для регистрации", "/menu")
	}

	link, err := GetLinkRegistrationByCode(slise[1])
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка параметра для регистрации", "/menu")
	}

	if !link.Active && link.Complete {
		webmessage.SendMessageClose(c, "данная ссылка уже была использована для регистрации", "/")
		return
	}

	if !link.Active {
		webmessage.SendMessageClose(c, "срок действия данной ссылки уже закончился", "/")
		return
	}

	var data struct {
		SchoolId         int
		SchoolName       string
		TypeRegistration string
		LinkCode         string
		ListCurrency     dict.ListCurrency
		Password         string
	}

	data.SchoolId = link.SchoolId
	data.TypeRegistration = "link"
	data.LinkCode = link.Code
	data.SchoolName = link.SchoolName

	data.ListCurrency = dict.GetAllCurrencyBySchoolId(c, data.SchoolId)

	data.Password = crypto.GeneratePassword()

	c.HTML(http.StatusOK, "i_teacher.html", data)

}

// обработка запроса на код
func GetDirectoryInputTeacherLinkHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	schoolId := c.Param("param")

	type Data struct {
		Link    string `json:"link"`
		Success bool   `json:"success"`
	}

	code := crypto.GetUlid()

	// Формируем ответ с Link и Success
	rs := Data{
		Link:    fmt.Sprintf("https://smart-crm.org.ua/input/teacher/registration/%s:%s", schoolId, code),
		Success: true, // Добавляем success
	}

	lr := LinkRegistration{}
	lr.SchoolId, _ = strconv.Atoi(schoolId)
	lr.Role = "teacher"
	lr.Code = code
	lr.Creator = user.Login

	lr.Rec(c)
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, rs)
}

// просмотр учителя
func GetDirectoryViewTeacherHandler(c *gin.Context) {

	teacher := dict.GetTeacherViewByLogin(c, internal.GetQueryString(c, "login"))
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "v_teacher.html", teacher)
}

// установить ставку учителя
func GetSetTeacherStavkaHandler(c *gin.Context) {

	lessonID := c.Query("lesson-id")
	stavka := c.Query("stavka")

	lessonIDInt, err := strconv.Atoi(lessonID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	stavkaInt, err := strconv.Atoi(stavka)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	err = dict.UpdateTeacherStavka(lessonIDInt, stavkaInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// установить ставку учителя
func GetSetTeacherPaimentHandler(c *gin.Context) {

	login := c.Query("login")
	piment := c.Query("piment")

	pimentInt, err := strconv.Atoi(piment)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	err = transaction.ExecTtransactionPaimentTeacher(c, login, pimentInt)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	teacher := dict.GetTeacherViewByLogin(c, login)
	if !c.IsAborted() {
		if teacher.TgId != 0 {
			mess := fmt.Sprintf(`
			<b>%s</b> <i>вам была произведена выплата</i> <b>%d %s</b>.`, teacher.FullName, pimentInt, teacher.CurrencyCod)
			tgbot.BotSendText(teacher.TgId, mess)
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
