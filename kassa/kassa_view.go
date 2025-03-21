package kassa

import (
	"encoding/json"
	"fmt"
	"schoolonline/postgree"

	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

type KassaView struct {
	ID             int       `db:"id" json:"id"`
	UserLogin      string    `db:"user_login" json:"user_login"`
	UserName       string    `db:"user_name" json:"user_name"`
	ConfirmedLogin string    `db:"confirmed_login_display" json:"confirmed_login"`
	ConfirmedName  string    `db:"confirmed_name" json:"confirmed_name"`
	IsDeleteLogin  string    `db:"is_delete_login_display" json:"is_delete_login"`
	IsDeleteName   string    `db:"is_delete_name" json:"is_delete_name"`
	SchoolId       int       `db:"school_id" json:"school_id"`
	SchoolName     string    `db:"school_name" json:"school_name"`
	CurrencyId     int       `db:"currency_id" json:"currency_id"`
	CurrencyCod    string    `db:"currency_cod" json:"currency_cod"`
	CurrencyName   string    `db:"currency_name" json:"currency_name"`
	Cost           int       `db:"cost" json:"cost"`
	Operation      string    `db:"operation" json:"operation"`
	OperationName  string    `db:"operation_name" json:"operation_name"`
	Comment        string    `db:"comment" json:"comment"`
	CreateTime     time.Time `db:"create_time" json:"create_time"`
	ConfirmeTime   time.Time `db:"confirme_time" json:"confirme_time"`
	IsDeleteTime   time.Time `db:"is_delete_time"`
	IsDeleted      bool      `db:"is_deleted" json:"is_deleted"`
}

type ListKassaView []KassaView

const queryReadKassaViewList = `
SELECT 
k.id,
k.user_login,
u.first_name || ' ' || u.last_name AS user_name,
COALESCE(NULLIF(TRIM(k.confirmed_login), ''), '-') AS confirmed_login_display,
COALESCE(c.first_name || ' ' || c.last_name, '-') AS confirmed_name, 
COALESCE(NULLIF(TRIM(k.is_delete_login), ''), '-') AS is_delete_login_display,
COALESCE(d.first_name || ' ' || d.last_name, '-') AS is_delete_name,
k.school_id,
s.name AS school_name,
k.currency_id,
k.currency_cod,
cr.name AS currency_name,
k.cost,
k.operation,
k.comment,
k.create_time,
k.confirme_time,
k.is_delete_time,
k.is_deleted
FROM kassa AS k
	LEFT JOIN school AS s ON s.id = k.school_id
	LEFT JOIN users AS u ON u.login = k.user_login
	LEFT JOIN users AS c ON c.login = k.confirmed_login
	LEFT JOIN users AS d ON d.login = k.is_delete_login
	LEFT JOIN currency AS cr ON cr.id = k.currency_id
`
const queryByIt = `
ORDER BY k.create_time DESC
LIMIT 1000
`

const queryByDirector = `
WHERE k.school_id = $1
ORDER BY k.create_time DESC
LIMIT 500
`

const queryByParent = `
WHERE u.login = $1
ORDER BY k.create_time DESC
LIMIT 50
`
const queryByTeacher = `
WHERE u.login = $1
ORDER BY k.create_time DESC
LIMIT 100
`

const queryByStudent = `
LEFT JOIN users AS p ON p.login = u.parent_login
WHERE p.login = $1
ORDER BY k.create_time DESC
LIMIT 50
`

func Get50KassaViewByRole(c *gin.Context, role, login string, schoolId int) (res ListKassaView) {

	var err error

	switch role {
	case "it":
		res, err = getKassaView(queryReadKassaViewList + queryByIt)
	case "director":
		res, err = getKassaView(queryReadKassaViewList+queryByDirector, fmt.Sprint(schoolId))
	case "parent":
		res, err = getKassaView(queryReadKassaViewList+queryByParent, login)
	case "teacher":
		res, err = getKassaView(queryReadKassaViewList+queryByTeacher, login)
	case "student":
		res, err = getKassaView(queryReadKassaViewList+queryByStudent, login)
	}
	if err != nil {
		webmessage.Err(c, err, "Ошибка при чтении кассы", "/directory")
		c.Abort()
		return
	}
	return
}

func getKassaView(query string, options ...string) (res ListKassaView, err error) {

	args := make([]interface{}, len(options))
	for i, v := range options {
		args[i] = v
	}

	// Выполняем запрос
	err = postgree.MainDBX.Select(&res, query, args...)
	// fmt.Println(options)
	// fmt.Println(args...)
	// fmt.Println(err)
	return
}

func (a *ListKassaView) SetOperationName() {
	if a == nil {
		return
	}
	for i, item := range *a {
		if operationName, ok := Operation[item.Operation]; ok {
			item.OperationName = item.Operation + " " + operationName
		} else {
			item.OperationName = item.Operation + " Unknown"
		}
		(*a)[i] = item
	}
}

func (a *ListKassaView) Marshall(c *gin.Context) string {
	r, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка json.Marshal pay", "/directory")
		c.Abort()
		return ""
	}
	return string(r)
}
