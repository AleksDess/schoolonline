package directory

import (
	"net/http"
	"schoolonline/dict"
	"schoolonline/internal"
	"schoolonline/routes/access"

	"github.com/gin-gonic/gin"
)

// список валют
func GetDirectoryListCurrencyHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	currency := dict.GetAllCurrencyByUser(c, user)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "l_currency.html", gin.H{
		"data": gin.H{"Currency": currency},
		"user": internal.GetUser(user),
	})
}

// добавление валюты
func PostDirectoryInputCurrencyHandler(c *gin.Context) {

	user := access.GetUser(c)

	sm := dict.GetIdSchoolByUser(c, user.Login)
	if c.IsAborted() {
		return
	}

	code := c.PostForm("code")
	symbol := c.PostForm("symbol")
	name := c.PostForm("name")
	rus_name := c.PostForm("rus_name")

	cur := dict.Currency{Code: code, SchoolId: sm.ID, Symbol: symbol, Name: name, RusName: rus_name}

	cur.Rec(c)
	if c.IsAborted() {
		return
	}

	// Возвращаем JSON-ответ
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Валюта добавлена",
	})
}

// просмотр валюты
func GetDirectoryViewCurrencyHandler(c *gin.Context) {

	// Получение currency по ID
	currency := dict.GetCurrencyByID(c)
	if c.IsAborted() {
		return
	}

	// Возврат HTML-контента в ответе
	c.HTML(http.StatusOK, "v_currency.html", currency)
}

// удаление currencyt
func GetDirectoryDeleteCurrencyHandler(c *gin.Context) {

	dict.DeleteCurrency(c)
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
