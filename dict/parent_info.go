package dict

import (
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"

	"github.com/gin-gonic/gin"
)

type ParentInfo struct {
	Login           string  `db:"login" json:"login"`
	FullName        string  `db:"full_name" json:"full_name"`
	SchoolName      string  `db:"school_name" json:"school_name"`
	SchoolId        int     `db:"school_id" json:"school_id"`
	CurrencyCod     string  `db:"currency_cod" json:"currency_cod"`
	CurrencyName    string  `db:"currency_name" json:"currency_name"`
	Balance         float64 `db:"balance" json:"balance"`
	TgId            int64   `db:"tg_id" json:"tg_id"`
	VbId            string  `db:"vb_id" json:"vb_id"`
	CountStudent    int     `db:"count_student" json:"count_student"`
	CountLesson     int     `db:"count_lesson" json:"count_lesson"`
	CountLessonWeek int     `db:"count_lesson_week" json:"count_lesson_week"`
	SummaLesson     int     `db:"summa_lesson" json:"summa_lesson"`
	SummaLessonWeek int     `db:"summa_lesson_week" json:"summa_lesson_week"`
}

const queryParentInfoByLogin = `
SELECT 
	u.login, 
	u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, 
	u.school_id, 
	u.balance, 
	c.code AS currency_cod,
	c.name AS currency_name,
	u.tg_id, 
	u.vb_id,
	(
        SELECT COUNT(*) 
        FROM users AS st 
        WHERE st.parent_login = u.login AND st.is_deleted = false
    ) AS count_student,
	(
        SELECT COUNT(*) 
        FROM lesson AS ls
		LEFT JOIN users AS us ON us.login = ls.student_login 
		LEFT JOIN users AS pr ON us.parent_login = us.login
        WHERE u.login = $1 AND ls.is_deleted = false
    ) AS count_lesson,
	(
        SELECT SUM(ls.count_lesson_per_week) 
        FROM lesson AS ls
		LEFT JOIN users AS us ON us.login = ls.student_login 
		LEFT JOIN users AS pr ON us.parent_login = us.login
        WHERE u.login = $1 AND ls.is_deleted = false
    ) AS count_lesson_week,
	(
        SELECT SUM(ls.cost_lesson) 
        FROM lesson AS ls
		LEFT JOIN users AS us ON us.login = ls.student_login 
		LEFT JOIN users AS pr ON us.parent_login = us.login
        WHERE u.login = $1 AND ls.is_deleted = false
    ) AS summa_lesson,
	(
        SELECT SUM(ls.cost_lesson) * SUM(ls.count_lesson_per_week)
        FROM lesson AS ls
		LEFT JOIN users AS us ON us.login = ls.student_login 
		LEFT JOIN users AS pr ON us.parent_login = us.login
        WHERE u.login = $1 AND ls.is_deleted = false
    ) AS summa_lesson_week
	FROM users AS u
		LEFT JOIN school AS s ON s.id = u.school_id
		LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.login = $1
`

func GetParentInfoByLogin(c *gin.Context, login string) (res ParentInfo) {

	err := postgree.MainDBX.Get(&res, queryParentInfoByLogin, login)

	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, nil, "Ошибка загрузки parentInfo из БД", "/directory")
		c.Abort()
		return
	}

	return
}

const queryParentInfoByTgId = `
SELECT 
	u.login, 
	u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, 
	u.school_id, 
	u.balance, 
	c.code AS currency_cod,
	c.name AS currency_name,
	u.tg_id, 
	u.vb_id,
	(
        SELECT COUNT(*) 
        FROM users AS st 
        WHERE st.parent_login = u.login AND st.is_deleted = false
    ) AS count_student,
	(
        SELECT COUNT(*) 
        FROM lesson AS ls
		LEFT JOIN users AS us ON us.login = ls.student_login 
		LEFT JOIN users AS pr ON us.parent_login = us.login
        WHERE u.tg_id = $1 AND ls.is_deleted = false
    ) AS count_lesson,
	(
        SELECT SUM(ls.count_lesson_per_week) 
        FROM lesson AS ls
		LEFT JOIN users AS us ON us.login = ls.student_login 
		LEFT JOIN users AS pr ON us.parent_login = us.login
        WHERE u.tg_id = $1 AND ls.is_deleted = false
    ) AS count_lesson_week,
	(
        SELECT SUM(ls.cost_lesson) 
        FROM lesson AS ls
		LEFT JOIN users AS us ON us.login = ls.student_login 
		LEFT JOIN users AS pr ON us.parent_login = us.login
        WHERE u.tg_id = $1 AND ls.is_deleted = false
    ) AS summa_lesson,
	(
        SELECT SUM(ls.cost_lesson) * SUM(ls.count_lesson_per_week)
        FROM lesson AS ls
		LEFT JOIN users AS us ON us.login = ls.student_login 
		LEFT JOIN users AS pr ON us.parent_login = us.login
        WHERE u.tg_id = $1 AND ls.is_deleted = false
    ) AS summa_lesson_week
	FROM users AS u
		LEFT JOIN school AS s ON s.id = u.school_id
		LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.tg_id = $1
`

func GetParentInfoByTgId(id int64) (res ParentInfo, err error) {

	err = postgree.MainDBX.Get(&res, queryParentInfoByTgId, id)

	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
