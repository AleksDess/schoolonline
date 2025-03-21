package menu

import (
	"net/http"
	"schoolonline/internal"
	"schoolonline/routes/access"

	"github.com/gin-gonic/gin"
)

func MenuGetHandler(c *gin.Context) {

	user := access.GetUser(c)

	c.HTML(http.StatusOK, "menu.html", gin.H{
		"user": internal.GetUser(user),
	})
}

type Link struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func GetMenuLinksHandler(c *gin.Context) {

	user := access.GetUser(c)

	links := []Link{}

	switch user.Role {
	case "it":
		links = LinksIt
	case "director":
		links = LinksDirector
	case "parent":
		links = LinksParent
	case "teacher":
		links = LinksTeacher
	case "student":
		links = LinksStudent
	}

	// Возвращаем данные в формате JSON
	c.JSON(http.StatusOK, links)
}

var LinksIt = []Link{
	{Text: "Справочники", URL: "/directory"},
	{Text: "Платежи", URL: "/pay"},
	{Text: "Уроки", URL: "/list/lesson"},
	{Text: "Выйти из аккаунта", URL: "/logout"},
}

var LinksDirector = []Link{
	{Text: "Справочники", URL: "/directory"},
	{Text: "Платежи", URL: "/pay"},
	{Text: "Уроки", URL: "/list/lesson"},
	{Text: "Выйти из аккаунта", URL: "/logout"},
}

var LinksParent = []Link{
	{Text: "Студенты", URL: "/list/student"},
	{Text: "Платежи", URL: "/pay"},
	{Text: "Выйти из аккаунта", URL: "/logout"},
}

var LinksTeacher = []Link{
	{Text: "Студенты", URL: "/list/student"},
	{Text: "Платежи", URL: "/pay"},
	{Text: "Уроки", URL: "/list/lesson"},
	{Text: "Выйти из аккаунта", URL: "/logout"},
}

var LinksStudent = []Link{
	{Text: "Уроки", URL: "/list/lesson"},
	{Text: "Выйти из аккаунта", URL: "/logout"},
}
