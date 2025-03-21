package internal

import (
	"encoding/json"
	"fmt"
	"schoolonline/dict"
	"schoolonline/webmessage"

	"github.com/gin-gonic/gin"
)

type User struct {
	Login    string `db:"login" json:"login" form:"login"`
	FullName string `db:"full_name" json:"full_name" form:"full_name"`
	Role     string `db:"role" json:"role" form:"role"`
}

func (a *User) Marshall(c *gin.Context) string {
	r, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка json.Marshal show", "/directory")
		c.Abort()
		return ""
	}
	return string(r)
}

func GetUser(u dict.User) *User {
	r := User{}
	r.Login = u.Login
	r.FullName = u.FirstName + " " + u.LastName
	r.Role = u.Role
	return &r
}

func GetUserJson(c *gin.Context, u dict.User) string {
	r := GetUser(u)
	return r.Marshall(c)
}
