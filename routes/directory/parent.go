package directory

import (
	"fmt"
	"net/http"
	"schoolonline/crypto"
	"schoolonline/dict"
	"schoolonline/internal"
	"schoolonline/routes/access"
	"schoolonline/transaction"
	"schoolonline/webmessage"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// список родителей
func GetDirectoryListParentHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	parents := dict.GetAllParentFromUser(c, user.Login)

	c.HTML(http.StatusOK, "l_parent.html", gin.H{
		"data":     parents,
		"user":     internal.GetUser(user),
		"schoolID": user.SchoolID,
	})
}

// удаление родителя
func GetDirectoryDeleteParentHandler(c *gin.Context) {

	dict.DeleteParent(c, internal.GetQueryString(c, "id"))
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ---------------------------------------------------
// добавление родителя пост запрос
func PostDirectoryInputParentHandler(c *gin.Context) {

	isUser := true

	user := access.GetUser(c)
	if c.IsAborted() {
		isUser = false
	}

	parent := dict.User{}

	if err := c.ShouldBind(&parent); err != nil {
		webmessage.Err(c, err, "Ошибка чтения формы", "/directory")
		c.Abort()
		return
	}

	us := dict.User{}
	us.Login = c.PostForm("login")
	us.PassWord = c.PostForm("password")
	us.Role = "parent"
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
	us.CreateTime = time.Now()
	us.UpdateTime = time.Now()
	us.VerifyEmail = true
	us.VerifyTime = time.Now()

	typeRegistration := c.PostForm("type_registration")
	linkCode := c.PostForm("link_code")

	switch typeRegistration {
	case "manual":
		transaction.ExecTransactionInsertUserAndParent(c, us)
		c.Redirect(http.StatusFound, "/list/parent")
	case "link":
		transaction.ExecTransactionInsertUserAndParentAndUpdateLink(c, us, linkCode)
		c.Redirect(http.StatusFound, "/menu")
	}
}

// добавление родителя вручную гет запрос
func GetDirectoryInputParentHandler(c *gin.Context) {

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

	c.HTML(http.StatusOK, "i_parent.html", data)
}

// добавление родителя по ссылке гет запрос
func GetDirectoryInputParentRegistration(c *gin.Context) {

	param := c.Param("param")

	slise := strings.Split(param, ":")

	if len(slise) != 2 {
		webmessage.Err(c, fmt.Errorf("invalid paraveters"), "неправильный параметра для регистрации", "/menu")
	}

	link, err := GetLinkRegistrationByCode(slise[1])
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка параметра для регистрации", "/menu")
		return
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

	c.HTML(http.StatusOK, "i_parent.html", data)
}

// обработка запроса на код
func GetDirectoryInputParentLinkHandler(c *gin.Context) {

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
		Link:    fmt.Sprintf("https://smart-crm.org.ua/input/parent/registration/%s:%s", schoolId, code),
		Success: true,
	}

	lr := LinkRegistration{}
	lr.SchoolId, _ = strconv.Atoi(schoolId)
	lr.Role = "parent"
	lr.Code = code
	lr.Creator = user.Login

	lr.Rec(c)
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, rs)
}

// просмотр родителя
func GetDirectoryViewParentHandler(c *gin.Context) {

	parent := dict.GetParentViewByLogin(c, internal.GetQueryString(c, "login"))
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "v_parent.html", parent)
}

// пополнения счета
func GetRefillParentBalanceHandler(c *gin.Context) {

	user := access.GetUser(c)

	parent := dict.GetParentViewByLogin(c, user.Login)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "refill_parent_balance.html", parent)
}
