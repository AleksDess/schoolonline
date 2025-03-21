package dict

import (
	"encoding/json"
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"

	"github.com/gin-gonic/gin"
)

type Student struct {
	Login          string  `db:"login" json:"login" form:"login"`
	FullName       string  `db:"full_name" json:"full_name" form:"full_name"`
	SchoolName     string  `db:"school_name" json:"school_name" form:"school_name"`
	SchoolId       int     `db:"school_id" json:"school_id" form:"school_id"`
	Email          string  `db:"email" json:"email" form:"email"`
	Phone          string  `db:"phone" json:"phone" form:"phone"`
	CurrencyCode   string  `db:"currency_code" json:"currency_code" form:"currency_code"`
	YearOfBirth    int     `db:"year_of_birth" json:"year_of_birth" form:"year_of_birth"`
	ParentLogin    string  `db:"parent_login" json:"parent_login" form:"parent_login"`
	ParentFullName string  `db:"parent_full_name" json:"parent_full_name" form:"parent_full_name" `
	ParentBalance  float64 `db:"parent_balance" json:"parent_balance" form:"parent_balance"`
	TgId           int64   `db:"tg_id" json:"tg_id" form:"tg_id"`
	VbId           string  `db:"vb_id" json:"vb_id" form:"vb_id"`
	LessonCost     int     `db:"lesson_cost" json:"lesson_cost" form:"lesson_cost"`
	LessonCount    int     `db:"lesson_count" json:"lesson_count" form:"lesson_count"`
}

type ListStudent []Student

func GetListStudentByUser(c *gin.Context, user string) (res ListStudent) {
	const query = `
		SELECT 
	u.login, 
	u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, u.school_id,
	u.email, u.phone, u.tg_id, u.vb_id,
	p.login AS parent_login,
	p.first_name || ' ' || p.last_name AS parent_full_name,
	c.code AS currency_code,
	p.balance AS parent_balance,
	u.year_of_birth,
	COALESCE(SUM(ls.cost_lesson), 0) AS lesson_cost,
	COALESCE(COUNT(ls.cost_lesson), 0) AS lesson_count
	FROM users AS u
	LEFT JOIN lesson AS ls ON ls.student_login = u.login
	LEFT JOIN users AS p ON p.login = u.parent_login
	LEFT JOIN school AS s ON s.id = u.school_id
	LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.school_id = (SELECT school_id FROM users WHERE login = $1)
	AND u.role = 'student'
	AND u.is_deleted = false
	AND ls.is_deleted = false
	AND ls.close = false
	GROUP BY 
    u.login, s.name, u.school_id, u.email, u.user_phone, u.tg_id, u.vb_id, 
    p.login, p.first_name, p.last_name, c.code, p.balance, u.year_of_birth;;`

	err := postgree.MainDBX.Select(&res, query, user)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка чтения student из БД", "/menu")
		c.Abort()
		return
	}
	return
}

func GetListStudentByStudentLogin(c *gin.Context, login string) (res ListStudent) {
	const query = `
		
		SELECT 
	u.login, 
	u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, u.school_id,
	u.email, u.phone, u.tg_id, u.vb_id,
	p.login AS parent_login,
	p.first_name || ' ' || p.last_name AS parent_full_name,
	c.code AS currency_code,
	p.balance AS parent_balance,
	u.year_of_birth,
	COALESCE(SUM(ls.cost_lesson), 0) AS lesson_cost,
	COALESCE(COUNT(ls.cost_lesson), 0) AS lesson_count
	FROM users AS u
	LEFT JOIN lesson AS ls ON ls.student_login = u.login
	LEFT JOIN users AS p ON p.login = u.parent_login
	LEFT JOIN school AS s ON s.id = u.school_id
	LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.login = $1
	AND u.role = 'student'
	AND u.is_deleted = false
	AND ls.is_deleted = false
	AND ls.close = false
	GROUP BY 
    u.login, s.name, u.school_id, u.email, u.user_phone, u.tg_id, u.vb_id, 
    p.login, p.first_name, p.last_name, c.code, p.balance, u.year_of_birth;;
	`

	err := postgree.MainDBX.Select(&res, query, login)
	if err != nil {

		fmt.Println(err)
		webmessage.Err(c, err, "ошибка чтения student из БД", "/menu")
		c.Abort()
		return
	}
	return
}

func GetListStudentByParent(c *gin.Context, parent string) (res ListStudent) {
	const query = `
		
		SELECT 
	u.login, 
	u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, u.school_id,
	u.email, u.phone, u.tg_id, u.vb_id,
	p.login AS parent_login,
	p.first_name || ' ' || p.last_name AS parent_full_name,
	c.code AS currency_code,
	p.balance AS parent_balance,
	u.year_of_birth,
	COALESCE(SUM(ls.cost_lesson), 0) AS lesson_cost,
	COALESCE(COUNT(ls.cost_lesson), 0) AS lesson_count
	FROM users AS u
	LEFT JOIN lesson AS ls ON ls.student_login = u.login
	LEFT JOIN users AS p ON p.login = u.parent_login
	LEFT JOIN school AS s ON s.id = u.school_id
	LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.parent_login = $1
	AND u.role = 'student'
	AND u.is_deleted = false
	GROUP BY 
    u.login, s.name, u.school_id, u.email, u.user_phone, u.tg_id, u.vb_id, 
    p.login, p.first_name, p.last_name, c.code, p.balance, u.year_of_birth;
	`

	err := postgree.MainDBX.Select(&res, query, parent)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка чтения student из БД", "/menu")
		c.Abort()
		return
	}
	return
}

func GetListStudentByTeacher(c *gin.Context, teacher string) (res ListStudent) {
	const query = `
		SELECT 
	u.login, 
	u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, u.school_id,
	u.email, u.phone, u.tg_id, u.vb_id,
	p.login AS parent_login,
	p.first_name || ' ' || p.last_name AS parent_full_name,
	c.code AS currency_code,
	p.balance AS parent_balance,
	u.year_of_birth,
	COALESCE(SUM(ls.cost_lesson), 0) AS lesson_cost,
	COALESCE(COUNT(ls.cost_lesson), 0) AS lesson_count
	FROM users AS u
	LEFT JOIN lesson AS ls ON ls.student_login = u.login
	LEFT JOIN users AS p ON p.login = u.parent_login
	LEFT JOIN school AS s ON s.id = u.school_id
	LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.school_id = (SELECT school_id FROM users WHERE login = $1)
	AND u.role = 'student'
	AND u.is_deleted = false
	AND ls.is_deleted = false
	AND ls.close = false
	GROUP BY 
    u.login, s.name, u.school_id, u.email, u.user_phone, u.tg_id, u.vb_id, 
    p.login, p.first_name, p.last_name, c.code, p.balance, u.year_of_birth;;`

	err := postgree.MainDBX.Select(&res, query, teacher)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка чтения student из БД", "/menu")
		c.Abort()
		return
	}
	return
}

func GetStudentByID(id int) (*Student, error) {
	const query = `
		SELECT 
	u.login, 
	u.first_name || ' ' || u.last_name AS full_name,
	s.name AS school_name, u.school_id,
	u.email, u.phone, u.tg_id, u.vb_id,
	p.login AS parent_login,
	p.first_name || ' ' || p.last_name AS parent_full_name,
	c.code AS currency_code,
	p.balance AS parent_balance,
	u.year_of_birth
	FROM users AS u
	LEFT JOIN users AS p ON p.login = u.parent_login
	LEFT JOIN school AS s ON s.id = u.school_id
	LEFT JOIN currency AS c ON c.id = u.currency_id
	WHERE u.login = $1;
	`

	var student Student
	err := postgree.MainDBX.Get(&student, query, id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &student, nil
}

// удалить student
func DeleteStudent(c *gin.Context, login string) {
	query := `UPDATE users SET is_deleted = true WHERE login = $1;`
	_, err := postgree.MainDB.Exec(query, login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка при удалении student", "/directory")
		c.Abort()
		return
	}
}

func GetListStudentByUserRole(c *gin.Context, user User) (stud ListStudent) {

	switch user.Role {
	case "it", "director":
		stud = GetListStudentByUser(c, user.Login)
		if c.IsAborted() {
			return
		}
	case "parent":
		stud = GetListStudentByParent(c, user.Login)
		if c.IsAborted() {
			return
		}
	case "teacher":
		stud = GetListStudentByTeacher(c, user.Login)
		if c.IsAborted() {
			return
		}
	default:
		webmessage.SendMessage(c, "У вас нет доступа к списку студентов", "/menu")
		return
	}
	return
}

func (a *ListStudent) Marshall(c *gin.Context) string {
	r, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка json.Marshal students", "/directory")
		c.Abort()
		return ""
	}
	return string(r)
}
