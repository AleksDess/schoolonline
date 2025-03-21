package dict

import (
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

type StudentView struct {
	Login           string         `db:"login" json:"login"`
	ParentLogin     string         `db:"parent_login" json:"parent_login"`
	SchoolID        int            `db:"school_id" json:"school_id"`
	YearOfBirth     int            `db:"year_of_birth" json:"year_of_birth"`
	CurrencyID      int            `db:"currency_id" json:"currency_id" form:"currency_id"`
	SchoolName      string         `db:"school_name" json:"school_name"`
	ParentFirstName string         `db:"parent_first_name" json:"parent_first_name"`
	ParentLastName  string         `db:"parent_last_name" json:"parent_last_name"`
	FirstName       string         `db:"first_name" json:"first_name"`
	LastName        string         `db:"last_name" json:"last_name"`
	CurrencyName    string         `db:"currency_name" json:"currency_name" form:"currency_name"`
	CreateTime      time.Time      `db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime      time.Time      `db:"update_time" json:"update_time" form:"update_time"`
	ListLesson      ListLessonView `db:"list_lesson" json:"list_lesson" form:"list_lesson"`
}

func GetStudentViewById(c *gin.Context, id string) (res StudentView) {
	const query = `
		SELECT u.login, u.parent_login, u.year_of_birth, u.first_name, u.last_name, u.create_time, u.update_time,
			p.first_name AS parent_first_name, p.last_name AS parent_last_name, p.currency_id, c.code AS currency_name,
			sc.id AS school_id, sc.name AS school_name
		FROM users AS u
		LEFT JOIN users AS p ON u.parent_login = p.login
		LEFT JOIN school AS sc ON p.school_id = sc.id
		LEFT JOIN currency AS c ON p.currency_id = c.id
		WHERE u.login = $1;`

	err := postgree.MainDBX.Get(&res, query, id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка загрузки student из БД", "/list/student")
		c.Abort()
		return
	}
	return
}
