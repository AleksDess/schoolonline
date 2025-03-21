package kassa

import (
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

var Operation = map[string]string{
	"1.1": "Зачисление от родителя",
	"1.2": "Списание у родителя",
	"3.1": "Зачисление учителю",
	"3.2": "Выплата учителю",
}

type Kassa struct {
	ID             int       `db:"id" json:"id"`
	UserLogin      string    `db:"user_login" json:"user_login"`
	ConfirmedLogin string    `db:"confirmed_login" json:"confirmed_login"`
	IsDeleteLogin  string    `db:"is_delete_login" json:"is_delete_login"`
	SchoolId       int       `db:"school_id" json:"school_id"`
	CurrencyId     int       `db:"currency_id" json:"currency_id"`
	CurrencyCod    string    `db:"currency_cod" json:"currency_cod"`
	Cost           int       `db:"cost" json:"cost"`
	Operation      string    `db:"operation" json:"operation"`
	Comment        string    `db:"comment" json:"comment"`
	CreateTime     time.Time `db:"create_time" json:"create_time"`
	ConfirmeTime   time.Time `db:"confirme_time" json:"confirme_time"`
	IsDeleteTime   time.Time `db:"is_delete_time"`
	IsDeleted      bool      `db:"is_deleted" json:"is_deleted"`
}

type ListKassa []Kassa

// функция создания таблицы
// постгресс User
func CreateTableKassa() error {

	const query = `
		CREATE TABLE IF NOT EXISTS kassa (
		id SERIAL PRIMARY KEY,
		user_login TEXT NOT NULL,
		confirmed_login TEXT NOT NULL,
		is_delete_login TEXT NOT NULL,
		school_id INTEGER NOT NULL,
		currency_id INTEGER NOT NULL,
		currency_cod TEXT NOT NULL,
		cost INTEGER NOT NULL,
		operation TEXT NOT NULL,
		comment TEXT NOT NULL,
		is_deleted BOOLEAN DEFAULT FALSE,
		create_time TIMESTAMPTZ,
		confirme_time TIMESTAMPTZ,
		is_delete_time TIMESTAMPTZ
		);
		`

	_, err := postgree.MainDBX.Exec(query)
	return err
}

const InsertKassaQuery = `
	INSERT INTO users
	(user_login, confirmed_login, is_delete_login, school_id, currency_id,
	cost, operation, comment, create_time, confirme_time, is_delete_time, is_deleted)
	VALUES 
	(:user_login, :confirmed_login, :is_delete_login, :school_id, :currency_id,
	:cost, :operation, :comment, :create_time, :confirme_time, :is_delete_time, :is_deleted)
`

func (a *Kassa) Rec(c *gin.Context) {

	a.CreateTime = time.Now()
	a.IsDeleteTime = time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)

	_, err := postgree.MainDBX.NamedExec(InsertKassaQuery, a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка записи kassa", "/directory")
		c.Abort()
		return
	}
}

func GetAllKassa(c *gin.Context) (res ListKassa) {
	err := postgree.MainDBX.Select(&res, "SELECT * FROM kassa;")
	if err != nil {
		webmessage.Err(c, err, "Ошибка при чтении кассы", "/directory")
		c.Abort()
		return
	}
	return
}

func GetKassaByID(id int) (res Kassa, err error) {
	err = postgree.MainDBX.Get(&res, "SELECT * FROM kassa WHERE id = $1;", id)
	return
}

const queryUpdateKassaIsDelete = `UPDATE lesson SET 
	is_deleted = true, 
	is_delete_time = $2 
	WHERE id = $1;`

func (a *Kassa) ExecIsDelete() (err error) {
	_, err = postgree.MainDBX.Exec(queryUpdateKassaIsDelete, a.ID, time.Now())
	return
}

const QuweryInsertKassaByTelegram = `
INSERT INTO kassa (
    user_login,
    confirmed_login,
    is_delete_login,
    school_id,
    currency_id,
	currency_cod,
    cost,
    operation,
	comment,
    create_time,
	confirme_time,
    is_delete_time,
    is_deleted
)
SELECT 
    u.login AS user_login,
    d.login AS confirmed_login, 
    ' ' AS is_delete_login, 
    u.school_id AS school_id,
    u.currency_id AS currency_id,
	c.code AS currency_cod,
    $3 AS cost,
    $4 AS operation,
	$5 AS comment,  
    $6 AS create_time,
	$7 AS confirme_time,
    $8 AS is_delete_time,
    $9 AS is_deleted
FROM 
    users AS u
	LEFT JOIN users AS d ON d.tg_id = $2
	LEFT JOIN currency AS c ON u.currency_id = c.id
WHERE 
    u.tg_id = $1; 
`
const QuweryInsertKassaByParentByStudentLogin = `
INSERT INTO kassa (
    user_login,
    confirmed_login,
    is_delete_login,
    school_id,
    currency_id,
	currency_cod,
    cost,
    operation,
	comment,
    create_time,
	confirme_time,
    is_delete_time,
    is_deleted
)
SELECT 
    p.login AS user_login,
    ' ' AS confirmed_login, 
    ' ' AS is_delete_login, 
    p.school_id AS school_id,
    p.currency_id AS currency_id,
	c.code AS currency_cod,
    $2 AS cost,
    $3 AS operation,
	$4 AS comment,  
    $5 AS create_time,
	$6 AS confirme_time,
    $7 AS is_delete_time,
    $8 AS is_deleted
FROM 
    users AS s
	LEFT JOIN users AS p ON s.parent_login = p.login
	LEFT JOIN currency AS c ON p.currency_id = c.id
WHERE 
    s.login = $1; 
`
const QuweryInsertKassaByTeacherByTeacherLogin = `
INSERT INTO kassa (
    user_login,
    confirmed_login,
    is_delete_login,
    school_id,
    currency_id,
	currency_cod,
    cost,
    operation,
	comment,
    create_time,
	confirme_time,
    is_delete_time,
    is_deleted
)
SELECT 
    t.login AS user_login,
    ' ' AS confirmed_login, 
    ' ' AS is_delete_login, 
    t.school_id AS school_id,
    t.currency_id AS currency_id,
	c.code AS currency_cod,
    $2 AS cost,
    $3 AS operation,
	$4 AS comment,  
    $5 AS create_time,
	$6 AS confirme_time,
    $7 AS is_delete_time,
    $8 AS is_deleted
FROM 
    users AS t
	LEFT JOIN currency AS c ON t.currency_id = c.id
WHERE 
    t.login = $1; 
`
