package webmessage

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Errors struct {
	Err          error
	ErrorMessage string
	RedirectURL  string
}

func Err(c *gin.Context, err error, mess, redir string) {
	e := Errors{}
	e.Err = err
	e.ErrorMessage = mess
	e.RedirectURL = redir
	WriteError(c, e)
}

func WriteError(c *gin.Context, r Errors) {

	if r.ErrorMessage == "" {
		r.ErrorMessage = "ошибка!!!"
	}

	if r.Err != nil {
		r.ErrorMessage += ":  (" + r.Err.Error() + ")"
	}

	if r.RedirectURL == "" {
		r.RedirectURL = "/menu"
	}

	c.HTML(http.StatusOK, "errors.html", r)
}

type Message struct {
	Message     string
	RedirectURL string
}

func SendMessage(c *gin.Context, mess, redir string) {
	e := Message{}
	e.Message = mess
	e.RedirectURL = redir
	WriteMessage(c, e)
}

func SendMessageClose(c *gin.Context, mess, redir string) {
	e := Message{}
	e.Message = mess
	e.RedirectURL = redir
	WriteMessageClose(c, e)
}

func WriteMessage(c *gin.Context, r Message) {

	if r.RedirectURL == "" {
		r.RedirectURL = "/menu"
	}

	c.HTML(http.StatusOK, "message.html", r)
}

func WriteMessageClose(c *gin.Context, r Message) {

	if r.RedirectURL == "" {
		r.RedirectURL = "/menu"
	}

	c.HTML(http.StatusOK, "messclose.html", r)
}

type MessageList struct {
	Titul    string
	Redirect string
	Message  []string
	Close    bool
	Err      error
}

func (a *MessageList) SendMessageList(c *gin.Context) {
	if a.Redirect == "" {
		a.Redirect = "/"
	}
	c.HTML(http.StatusOK, "messageList.html", a)
}
