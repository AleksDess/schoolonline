package directory

import (
	"net/http"
	"schoolonline/dict"
	"schoolonline/internal"
	"schoolonline/routes/access"
	"schoolonline/webmessage"

	"github.com/gin-gonic/gin"
)

// список предмет
func GetDirectoryListItemHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	items := dict.GetAllItemsFromUser(c, user.Login)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "l_item.html", gin.H{
		"Items": items,
		"user":  internal.GetUser(user),
	})
}

// удаление предмета
func GetDirectoryDeleteItemHandler(c *gin.Context) {

	dict.DeleteItem(c)
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// добавление предмета
func GetDirectoryInputItemHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	res := dict.GetListFacultySmSchoolSmByUser(c, user.Login)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "i_item.html", gin.H{
		"jsonData": string(res.Marshall(c)),
	})
}

// добавление предмета
func PostDirectoryInputItemHandler(c *gin.Context) {

	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}

	data := dict.Item{}

	if err := c.ShouldBind(&data); err != nil {
		webmessage.Err(c, err, "Ошибка чтения формы", "/list/item")
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

	c.Redirect(http.StatusFound, "/list/item")
}

// просмотр предмета
func GetDirectoryViewItemHandler(c *gin.Context) {

	item := dict.GetItemByID(c)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, "v_item.html", item)
}
