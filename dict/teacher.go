package dict

import (
	"encoding/json"
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

type Teacher struct {
	Login         string  `db:"login" json:"login" form:"login"`
	FullName      string  `db:"full_name" json:"full_name" form:"full_name"`
	SchoolName    string  `db:"school_name" json:"school_name" form:"school_name"`
	SchoolId      int     `db:"school_id" json:"school_id" form:"school_id"`
	Email         string  `db:"email" json:"email" form:"email"`
	Phone         string  `db:"phone" json:"phone" form:"phone"`
	CurrencyName  string  `db:"currency_name" json:"currency_name" form:"currency_name"`
	CurrencyCod   string  `db:"currency_cod" json:"currency_cod" form:"currency_cod"`
	PimentDetails string  `db:"piment_details" json:"piment_details" form:"piment_details"`
	Balance       float64 `db:"balance" json:"balance" form:"balance"`
	TgId          int64   `db:"tg_id" json:"tg_id" form:"tg_id"`
	VbId          string  `db:"vb_id" json:"vb_id" form:"vb_id"`
}

type ListTeacher []Teacher

// GetAllTeachersBySchoolId
func GetAllTeachersBySchoolId(c *gin.Context, id string) (res ListTeacher) {
	const query = `
		SELECT u.login, u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, u.school_id, u.email, u.phone, u.tg_id, u.vb_id,
	c.code AS currency_cod, c.name AS currency_name, u.piment_details, u.balance
	FROM users AS u
	LEFT JOIN school AS s ON s.id = u.school_id
	LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.school_id = $1
	AND u.role = 'teacher'
	AND u.is_deleted = false;
		`

	err := postgree.MainDBX.Select(&res, query, id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка чтения учителей школы", "list/student")
		c.Abort()
		return
	}
	return
}

func GetAllTeacherFromUser(c *gin.Context, user string) (res ListTeacher) {
	const query = `
		SELECT u.login, u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, u.school_id, u.email, u.phone, u.tg_id, u.vb_id,
	c.code AS currency_cod, c.name AS currency_name, u.piment_details, u.balance
	FROM users AS u
	LEFT JOIN school AS s ON s.id = u.school_id
	LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.school_id = (SELECT school_id FROM users WHERE login = $1)
	AND u.role = 'teacher'
	AND u.is_deleted = false;
	`

	err := postgree.MainDBX.Select(&res, query, user)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка списка учителей школы", "/directory")
		c.Abort()
		return
	}
	return
}

func GetTeacherViewByLogin(c *gin.Context, id string) (res Teacher) {
	const query = `
		SELECT u.login, u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, u.school_id, u.email, u.phone, 
	u.balance, c.code AS currency_cod, c.name AS currency_name, u.tg_id, u.vb_id, u.piment_details
	FROM users AS u
	LEFT JOIN school AS s ON s.id = u.school_id
	LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.login = $1;
	`

	err := postgree.MainDBX.Get(&res, query, id)

	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, nil, "Ошибка загрузки учителя из БД", "/directory")
		c.Abort()
		return
	}

	return
}

func GetTeacherByLogin(c *gin.Context, login string) (res User) {
	const query = `
		SELECT * FROM users WHERE login = $1;`

	err := postgree.MainDBX.Get(&res, query, login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка загрузки учителя из БД", "/directory")
		c.Abort()
		return
	}
	return
}

// удалить родителя
func DeleteTeacher(c *gin.Context, login string) {
	query := `UPDATE users SET is_deleted = true, update_time = $2 WHERE login = $1;`
	_, err := postgree.MainDB.Exec(query, login, time.Now())
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка при удалении учителя", "/directory")
		c.Abort()
		return
	}
}

func (a *ListTeacher) Marshall(c *gin.Context) string {
	r, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка json.Marshal lesson view", "/directory")
		c.Abort()
		return ""
	}
	return string(r)
}
