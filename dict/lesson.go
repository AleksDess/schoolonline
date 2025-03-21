package dict

import (
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// метод учета - инт
//               0  -  не определено
//               1  -  процент учителя
//				 2  -  ставка учителя

type Lesson struct {
	ID                 int       `db:"id" json:"id" form:"id"`
	ItemID             int       `db:"item_id" json:"item_id" form:"item_id"`
	StudentLogin       string    `db:"student_login" json:"student_login" form:"student_login"`
	TeacherLogin       string    `db:"teacher_login" json:"teacher_login" form:"teacher_login"`
	CostLesson         int       `db:"cost_lesson" json:"cost_lesson" form:"cost_lesson"`
	CountLessonPerWeek int       `db:"count_lesson_per_week" json:"count_lesson_per_week" form:"count_lesson_per_week"`
	TeacherBaseRate    int       `db:"teacher_base_rate" json:"teacher_base_rate" form:"teacher_base_rate"`
	IsGroup            bool      `db:"is_group" json:"is_group" form:"is_group"`
	GroupID            int       `db:"group_id" json:"group_id" form:"group_id"`
	Close              bool      `db:"close" json:"close" form:"close"`
	IsDeleted          bool      `db:"is_deleted" json:"is_deleted" form:"is_deleted"`
	CreateTime         time.Time `db:"create_time" json:"create_time" form:"create_time"`
	CloseTime          time.Time `db:"close_time" json:"close_time" form:"close_time"`
	StartDate          time.Time `db:"start_date" json:"start_date" form:"start_date" time_format:"2006-01-02"`
	EndDate            time.Time `db:"end_date" json:"end_date" form:"end_date" time_format:"2006-01-02"`
	DurationMinutes    int       `db:"duration_minutes" json:"duration_minutes" form:"duration_minutes"`
	Description        string    `db:"description" json:"description" form:"description"`
}

type ListLesson []Lesson

func (lesson *Lesson) Print() {
	fmt.Printf("ID:                  %d\n", lesson.ID)
	fmt.Printf("ItemID:              %d\n", lesson.ItemID)
	fmt.Printf("StudentLogin:        %s\n", lesson.StudentLogin)
	fmt.Printf("TeacherLogin:        %s\n", lesson.TeacherLogin)
	fmt.Printf("CostLesson:          %d\n", lesson.CostLesson)
	fmt.Printf("TeacherBaseRate:     %d\n", lesson.TeacherBaseRate)
	fmt.Printf("IsGroup:             %t\n", lesson.IsGroup)
	fmt.Printf("GroupID:             %d\n", lesson.GroupID)
	fmt.Printf("Close:               %t\n", lesson.Close)
	fmt.Printf("IsDeleted:           %t\n", lesson.IsDeleted)
	fmt.Printf("CreateTime:          %s\n", lesson.CreateTime)
	fmt.Printf("CloseTime:           %s\n", lesson.CloseTime)
	fmt.Printf("StartDate:           %s\n", lesson.StartDate)
	fmt.Printf("EndDate:             %s\n", lesson.EndDate)
	fmt.Printf("DurationMinutes:     %d\n", lesson.DurationMinutes)
	fmt.Printf("Description:         %s\n", lesson.Description)
}

func CreateLessonTableIfNotExists() error {
	const query = `
        CREATE TABLE IF NOT EXISTS lesson (
            id SERIAL PRIMARY KEY,
            item_id INT NOT NULL,
            student_login TEXT NOT NULL,
            teacher_login TEXT NOT NULL,
            cost_lesson INT NOT NULL,
			count_lesson_per_week INT NOT NULL,
            teacher_base_rate INT NOT NULL,
            is_group BOOLEAN NOT NULL,
            group_id INT NOT NULL,
            close BOOLEAN NOT NULL,
            is_deleted BOOLEAN NOT NULL,
            create_time TIMESTAMP NOT NULL,
            close_time TIMESTAMP NOT NULL,
            start_date TIMESTAMP NOT NULL,
            end_date TIMESTAMP,
            duration_minutes INT NOT NULL,
            description TEXT
        );
    `

	_, err := postgree.MainDBX.Exec(query)
	return err
}

func GetAllLesson() (res ListLesson, err error) {
	err = postgree.MainDBX.Select(&res, "SELECT * FROM lesson;")
	return
}

func GetAllLessonByStudentLogin(login string) (res ListLesson, err error) {
	err = postgree.MainDBX.Select(&res, "SELECT * FROM lesson WHERE student_login = $1;", login)
	return
}

func GetLessonById(id int) (res Lesson, err error) {
	err = postgree.MainDBX.Get(&res, "SELECT * FROM lesson WHERE id = $1;", id)
	return
}

func UpdateTeacherStavka(id, cost int) error {
	_, err := postgree.MainDB.Exec("UPDATE lesson SET teacher_base_rate = $2 WHERE id = $1", id, cost)
	return err
}

func GetLessonByIdWeb(c *gin.Context) (res Lesson) {

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

	err := postgree.MainDBX.Get(&res, "SELECT * FROM lesson WHERE id = $1;", data)
	if err != nil {
		fmt.Println(err)
		c.Abort()
		return
	}
	return
}

const queryInsertLesson = `
        INSERT INTO lesson (
            item_id, student_login, teacher_login, cost_lesson, count_lesson_per_week, teacher_base_rate, 
			is_group, group_id, close, is_deleted, create_time, close_time, start_date, end_date, duration_minutes, description
        )
        VALUES (
            :item_id, :student_login, :teacher_login, :cost_lesson, :count_lesson_per_week, :teacher_base_rate,
			:is_group, :group_id, :close, :is_deleted, :create_time, :close_time, :start_date, :end_date, :duration_minutes, :description
        );
`

func (lesson *Lesson) Rec(c *gin.Context) {

	_, err := postgree.MainDBX.NamedExec(queryInsertLesson, lesson)

	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка записи урока", "/list/student")
		c.Abort()
		return
	}
}

func (a *Lesson) CheckLessonCompletionDate() bool {

	query := `SELECT COUNT(*) 
FROM lesson_complete 
WHERE lesson_id = $1 
  AND DATE(create_time) = DATE($2) 
  AND is_deleted = FALSE;`

	query = strings.Replace(query, "$1", "'"+fmt.Sprint(a.ID)+"'", 1)
	query = strings.Replace(query, "$2", "'"+fmt.Sprint(time.Now().Format("2006-01-02 15:05:04"))+"'", 1)

	var n int
	row := postgree.MainDB.QueryRow(query)
	err := row.Scan(&n)
	if err != nil {
		fmt.Println("***", err)
		return true
	}

	if n > 0 {
		return true
	}
	return false
}
