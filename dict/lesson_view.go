package dict

import (
	"encoding/json"
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

type LessonView struct {
	ID                  int       `db:"id" json:"id" form:"id"`
	ItemId              int       `db:"item_id" json:"item_id" form:"item_id"`
	TeacherLogin        string    `db:"teacher_login" json:"teacher_login" form:"teacher_login"`
	ParentLogin         string    `db:"parent_login" json:"parent_login" form:"parent_login"`
	ItemName            string    `db:"item_name" json:"item_name" form:"item_name"`
	ParentName          string    `db:"parent_name" json:"parent_name" form:"parent_name"`
	StudentName         string    `db:"student_name" json:"student_name" form:"student_name"`
	TeacherName         string    `db:"teacher_name" json:"teacher_name" form:"teacher_name"`
	CurrencyName        string    `db:"currency_name" json:"currency_name" form:"currency_name"`
	CurrencyNameTeacher string    `db:"currency_name_teacher" json:"currency_name_teacher" form:"currency_name_teacher"`
	CostLesson          int       `db:"cost_lesson" json:"cost_lesson" form:"cost_lesson"`
	CountLessonPerWeek  int       `db:"count_lesson_per_week" json:"count_lesson_per_week" form:"count_lesson_per_week"`
	TeacherBaseRate     int       `db:"teacher_base_rate" json:"teacher_base_rate" form:"teacher_base_rate"`
	ParentBalance       float64   `db:"parent_balance" json:"parent_balance" form:"parent_balance"`
	IsGroup             bool      `db:"is_group" json:"is_group" form:"is_group"`
	Close               bool      `db:"close" json:"close" form:"close"`
	IsDeleted           bool      `db:"is_deleted" json:"is_deleted" form:"is_deleted"`
	CreateTime          time.Time `db:"create_time" json:"create_time" form:"create_time"`
	StartDate           time.Time `db:"start_date" json:"start_date" form:"start_date"`
	EndDate             time.Time `db:"end_date" json:"end_date" form:"end_date"`
	DurationMinutes     int       `db:"duration_minutes" json:"duration_minutes" form:"duration_minutes"`
	Description         string    `db:"description" json:"description" form:"description"`
}

type ListLessonView []LessonView

func GetLessonViewById(c *gin.Context) (res LessonView) {

	id := c.Query("id")
	login := c.Query("login")

	if id == "" && login == "" {
		webmessage.Err(c, fmt.Errorf("id not provided"), "ID урока не указан", "/directory")
		c.Abort()
		return
	}

	data := ""

	if id == "" {
		data = login
	} else {
		data = id
	}

	const query = `
		SELECT 
			l.id,
			l.item_id,
			l.teacher_login,
			s.parent_login,
			i.name AS item_name,
			t.first_name || ' ' || t.last_name AS teacher_name,
			p.first_name || ' ' || p.last_name AS parent_name,
			s.first_name || ' ' || s.last_name AS student_name,
			c.code AS currency_name,
			ct.code AS currency_name_teacher,
			l.cost_lesson,
			l.count_lesson_per_week,
			l.duration_minutes,
			l.teacher_base_rate,
			p.balance AS parent_balance,
			l.is_group,
			l.close,
			l.is_deleted,
			l.create_time,
			l.start_date,
			l.end_date,
			l.duration_minutes,
			l.description
		FROM lesson AS l
		LEFT JOIN item AS i ON l.item_id = i.id
		LEFT JOIN users AS s ON l.student_login = s.login
		LEFT JOIN users AS t ON l.teacher_login = t.login
		LEFT JOIN users AS p ON s.parent_login = p.login
		LEFT JOIN currency AS c ON p.currency_id = c.id
		LEFT JOIN currency AS ct ON t.currency_id = ct.id
		WHERE l.id = $1;
		`

	err := postgree.MainDBX.Get(&res, query, data)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка загрузки Lesson из БД", "/list/lesson")
		c.Abort()
		return
	}
	return
}

func GetAllLessonViewByStudentLogin(c *gin.Context, login string) (res ListLessonView) {
	const query = `
		SELECT 
			l.id,
			l.item_id,
			l.teacher_login,
			i.name AS item_name,
			t.first_name || ' ' || t.last_name AS teacher_name,
			p.first_name || ' ' || p.last_name AS parent_name,
			s.first_name || ' ' || s.last_name AS student_name,
			c.code AS currency_name,
			ct.code AS currency_name_teacher,
			l.cost_lesson,
			l.count_lesson_per_week,
			l.duration_minutes,
			l.teacher_base_rate,
			p.balance AS parent_balance,
			l.is_group,
			l.close,
			l.is_deleted,
			l.create_time,
			l.start_date,
			l.end_date,
			l.duration_minutes,
			l.description
		FROM lesson AS l
		LEFT JOIN item AS i ON l.item_id = i.id
		LEFT JOIN users AS s ON l.student_login = s.login
		LEFT JOIN users AS t ON l.teacher_login = t.login
		LEFT JOIN users AS p ON s.parent_login = p.login
		LEFT JOIN currency AS c ON p.currency_id = c.id
		LEFT JOIN currency AS ct ON t.currency_id = ct.id
		WHERE l.student_login = $1
			AND l.close = false
			AND l.is_deleted = false;`

	err := postgree.MainDBX.Select(&res, query, login)

	if err != nil {

		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка загрузки listLesson из БД", "/list/student")
		c.Abort()
		return
	}
	return
}

func GetAllLessonViewByUserLoginDirector(c *gin.Context, login string) (res ListLessonView) {
	const query = `
		SELECT 
			l.id,
			l.item_id,
			l.teacher_login,
			i.name AS item_name,
			t.first_name || ' ' || t.last_name AS teacher_name,
			p.first_name || ' ' || p.last_name AS parent_name,
			s.first_name || ' ' || s.last_name AS student_name,
			c.code AS currency_name,
			ct.code AS currency_name_teacher,
			l.cost_lesson,
			l.count_lesson_per_week,
			l.duration_minutes,
			l.teacher_base_rate,
			p.balance AS parent_balance,
			l.is_group,
			l.close,
			l.is_deleted,
			l.create_time,
			l.start_date,
			l.end_date,
			l.duration_minutes,
			l.description
		FROM lesson AS l
		LEFT JOIN item AS i ON l.item_id = i.id
		LEFT JOIN users AS t ON l.teacher_login = t.login
		LEFT JOIN users AS s ON l.student_login = s.login
		LEFT JOIN users AS p ON s.parent_login = p.login
		LEFT JOIN currency AS c ON p.currency_id = c.id
		LEFT JOIN currency AS ct ON t.currency_id = ct.id
		LEFT JOIN school AS sc ON s.school_id = sc.id
		LEFT JOIN users AS u ON u.login = sc.director
		WHERE u.login = $1
			AND l.close = false
			AND l.is_deleted = false;`

	err := postgree.MainDBX.Select(&res, query, login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка загрузки listLesson из БД", "/directory")
		c.Abort()
		return
	}
	return
}

func GetAllLessonViewByUserLoginParent(c *gin.Context, login string) (res ListLessonView) {
	const query = `
		SELECT 
			l.id,
			l.item_id,
			l.teacher_login,
			i.name AS item_name,
			t.first_name || ' ' || t.last_name AS teacher_name,
			p.first_name || ' ' || p.last_name AS parent_name,
			s.first_name || ' ' || s.last_name AS student_name,
			c.code AS currency_name,
			ct.code AS currency_name_teacher,
			l.cost_lesson,
			l.count_lesson_per_week,
			l.duration_minutes,
			l.teacher_base_rate,
			p.balance AS parent_balance,
			l.is_group,
			l.close,
			l.is_deleted,
			l.create_time,
			l.start_date,
			l.end_date,
			l.duration_minutes,
			l.description
		FROM lesson AS l
		LEFT JOIN item AS i ON l.item_id = i.id
		LEFT JOIN users AS t ON l.teacher_login = t.login
		LEFT JOIN users AS s ON l.student_login = s.login
		LEFT JOIN users AS p ON s.parent_login = p.login
		LEFT JOIN currency AS c ON p.currency_id = c.id
		LEFT JOIN currency AS ct ON t.currency_id = ct.id
		LEFT JOIN school AS sc ON s.school_id = sc.id
		LEFT JOIN users AS u ON u.login = sc.director
		WHERE p.login = $1
			AND l.close = false
			AND l.is_deleted = false;`

	err := postgree.MainDBX.Select(&res, query, login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка загрузки listLesson из БД", "/directory")
		c.Abort()
		return
	}
	return
}

func GetAllLessonViewByUserLoginTeacher(c *gin.Context, login string) (res ListLessonView) {
	const query = `
		SELECT 
			l.id,
			l.item_id,
			l.teacher_login,
			i.name AS item_name,
			t.first_name || ' ' || t.last_name AS teacher_name,
			p.first_name || ' ' || p.last_name AS parent_name,
			s.first_name || ' ' || s.last_name AS student_name,
			c.code AS currency_name,
			ct.code AS currency_name_teacher,
			l.cost_lesson,
			l.count_lesson_per_week,
			l.duration_minutes,
			l.teacher_base_rate,
			p.balance AS parent_balance,
			l.is_group,
			l.close,
			l.is_deleted,
			l.create_time,
			l.start_date,
			l.end_date,
			l.duration_minutes,
			l.description
		FROM lesson AS l
		LEFT JOIN item AS i ON l.item_id = i.id
		LEFT JOIN users AS t ON l.teacher_login = t.login
		LEFT JOIN users AS s ON l.student_login = s.login
		LEFT JOIN users AS p ON s.parent_login = p.login
		LEFT JOIN currency AS c ON p.currency_id = c.id
		LEFT JOIN currency AS ct ON t.currency_id = ct.id
		LEFT JOIN school AS sc ON s.school_id = sc.id
		LEFT JOIN users AS u ON u.login = sc.director
		WHERE t.login =  $1
			AND l.close = false
			AND l.is_deleted = false;`

	err := postgree.MainDBX.Select(&res, query, login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка загрузки listLesson из БД", "/directory")
		c.Abort()
		return
	}
	return
}

func GetAllLessonViewByUserLoginIt(c *gin.Context) (res ListLessonView) {
	const query = `
		SELECT 
			l.id,
			l.item_id,
			l.teacher_login,
			i.name AS item_name,
			t.first_name || ' ' || t.last_name AS teacher_name,
			p.first_name || ' ' || p.last_name AS parent_name,
			s.first_name || ' ' || s.last_name AS student_name,
			c.code AS currency_name,
			ct.code AS currency_name_teacher,
			l.cost_lesson,
			l.count_lesson_per_week,
			l.duration_minutes,
			l.teacher_base_rate,
			p.balance AS parent_balance,
			l.is_group,
			l.close,
			l.is_deleted,
			l.create_time,
			l.start_date,
			l.end_date,
			l.duration_minutes,
			l.description
		FROM lesson AS l
		LEFT JOIN item AS i ON l.item_id = i.id
		LEFT JOIN users AS t ON l.teacher_login = t.login
		LEFT JOIN users AS s ON l.student_login = s.login
		LEFT JOIN users AS p ON s.parent_login = p.login
		LEFT JOIN currency AS c ON p.currency_id = c.id
		LEFT JOIN currency AS ct ON t.currency_id = ct.id
		WHERE  l.close = false
			AND l.is_deleted = false;`

	err := postgree.MainDBX.Select(&res, query)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка загрузки listLesson из БД", "/directory")
		c.Abort()
		return
	}
	return
}

// закрыть урок
func UpdateLessonClose(c *gin.Context) string {

	id := c.Query("id")
	if id == "" {
		webmessage.Err(c, fmt.Errorf("id not provided"), "ID предмета не указан", "/directory")
		c.Abort()
		return ""
	}

	query := `UPDATE lesson	SET close = true, end_date = $2 WHERE id = $1;`

	_, err := postgree.MainDBX.Exec(query, id, time.Now())
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка обновления clise is Lesson", "/directory")
		c.Abort()
		return ""
	}

	return id
}

func (a *ListLessonView) Marshall(c *gin.Context) string {
	r, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка json.Marshal lesson view", "/directory")
		c.Abort()
		return ""
	}
	return string(r)
}

func (a *LessonView) Marshall(c *gin.Context) string {
	r, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка json.Marshal lesson view", "/directory")
		c.Abort()
		return ""
	}
	return string(r)
}
