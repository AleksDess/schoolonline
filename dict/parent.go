package dict

import (
	"encoding/json"
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"

	"github.com/gin-gonic/gin"
)

type Parent struct {
	Login       string  `db:"login" json:"login" form:"login"`
	FullName    string  `db:"full_name" json:"full_name" form:"full_name"`
	SchoolName  string  `db:"school_name" json:"school_name" form:"school_name"`
	SchoolId    int     `db:"school_id" json:"school_id" form:"school_id"`
	Email       string  `db:"email" json:"email" form:"email"`
	Phone       string  `db:"phone" json:"phone" form:"phone"`
	CurrencyCod string  `db:"currency_cod" json:"currency_cod" form:"currency_cod"`
	Balance     float64 `db:"balance" json:"balance" form:"balance"`
	TgId        int64   `db:"tg_id" json:"tg_id" form:"tg_id"`
	VbId        string  `db:"vb_id" json:"vb_id" form:"vb_id"`
}

type ListParent []Parent

func GetAllParentFromUser(c *gin.Context, user string) ListParent {
	const query = `SELECT 
	u.login, 
	u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, 
	u.school_id, 
	u.email, 
	u.phone, 
	u.tg_id, 
	u.vb_id
	FROM users AS u
		LEFT JOIN school AS s ON s.id = u.school_id
	WHERE u.school_id = (SELECT school_id FROM users WHERE login = $1)
	AND u.role = 'parent'
	AND u.is_deleted = false;
	`

	var res ListParent
	err := postgree.MainDBX.Select(&res, query, user)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка списка родителей", "/directory")
		c.Abort()
		return nil
	}
	return res
}

func GetParentFromUser(c *gin.Context, user string) ListParent {
	const query = `SELECT 
	u.login, 
	u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, 
	u.school_id, 
	u.email, 
	u.phone, 
	u.tg_id, 
	u.vb_id
	FROM users AS u
		LEFT JOIN school AS s ON s.id = u.school_id
	WHERE u.login = $1
	AND u.role = 'parent'
	AND u.is_deleted = false;
	`

	var res ListParent
	err := postgree.MainDBX.Select(&res, query, user)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка списка родителей", "/directory")
		c.Abort()
		return nil
	}
	return res
}

func (a *ListParent) Marshall(c *gin.Context) string {
	r, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка json.Marshal школ и предметов", "/directory")
		c.Abort()
		return ""
	}
	return string(r)
}

// удалить родителя
func DeleteParent(c *gin.Context, login string) {

	if login == "" {
		webmessage.Err(c, fmt.Errorf("error"), "Ошибка при удалении родителя", "/directory")
		c.Abort()
		return
	}

	query := `UPDATE users SET is_deleted = true 
	WHERE login = $1
	OR parent_login = $1;`
	_, err := postgree.MainDB.Exec(query, login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка при удалении родителя", "/directory")
		c.Abort()
		return
	}
}

func GetParentViewByLogin(c *gin.Context, id string) (res Parent) {
	const query = `SELECT 
	u.login, 
	u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, 
	u.school_id, 
	u.email, 
	u.phone, 
	u.balance, 
	c.code AS currency_cod, 
	u.tg_id, 
	u.vb_id
	FROM users AS u
		LEFT JOIN school AS s ON s.id = u.school_id
		LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.login = $1;
	`

	err := postgree.MainDBX.Get(&res, query, id)

	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, nil, "Ошибка загрузки родителя из БД", "/directory")
		c.Abort()
		return
	}

	return
}

type ParentInfoFromAccauntReplenishment struct {
	Name            string `db:"name"`
	CurrencyRusName string `db:"rus_name"`
	CurrencyCode    string `db:"code"`
}

// для пополнения счета
func GetParentNameAndCurrencyNameByLogin(login string) (res ParentInfoFromAccauntReplenishment, err error) {
	const query = `
		SELECT  p.first_name || ' ' || p.last_name AS name,
			c.rus_name,
			c.code
		FROM users AS p
			LEFT JOIN currency AS c ON c.id = p.currency_id
		WHERE p.login = $1;
	`
	err = postgree.MainDBX.Get(&res, query, login)
	return
}
