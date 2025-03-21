package directory

import (
	"net/http"
	"schoolonline/internal"
	"schoolonline/routes/access"

	"github.com/gin-gonic/gin"
)

func GetDirectoryHandler(c *gin.Context) {

	user := access.GetUser(c)

	c.HTML(http.StatusOK, "dict.html", gin.H{
		"user": internal.GetUser(user),
	})
}

type Link struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func GetDirectoryLinksHandler(c *gin.Context) {

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
	}

	// Возвращаем данные в формате JSON
	c.JSON(http.StatusOK, links)
}

var LinksIt = []Link{
	{Text: "Школа", URL: "/list/school"},
	{Text: "Факультет", URL: "/list/faculty"},
	{Text: "Предмет", URL: "/list/item"},
	{Text: "Родители", URL: "/list/parent"},
	{Text: "Студенты", URL: "/list/student"},
	{Text: "Учителя", URL: "/list/teacher"},
	{Text: "Валюта", URL: "/list/currency"},
}

var LinksDirector = []Link{
	{Text: "Школа", URL: "/list/school"},
	{Text: "Факультет", URL: "/list/faculty"},
	{Text: "Предмет", URL: "/list/item"},
	{Text: "Родители", URL: "/list/parent"},
	{Text: "Студенты", URL: "/list/student"},
	{Text: "Учителя", URL: "/list/teacher"},
	{Text: "Валюта", URL: "/list/currency"},
}

var LinksParent = []Link{
	{Text: "Студенты", URL: "/list/student"},
}

var LinksTeacher = []Link{
	{Text: "Факультет", URL: "/list/faculty"},
	{Text: "Предмет", URL: "/list/item"},
}
