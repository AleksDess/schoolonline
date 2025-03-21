package dict

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Faculty struct {
	ID          int       `db:"id" json:"id"`
	School      int       `form:"school" db:"school" json:"school"`
	Name        string    `form:"name" db:"name" json:"name"`
	Description string    `form:"description" db:"description" json:"description"`
	IsDeleted   bool      `db:"is_deleted" json:"is_deleted" form:"is_deleted"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	SchoolName  string    `db:"school_name"`
}

type ListFaculty []Faculty

func CreateFacultyTableIfNotExists() error {
	const query = `
		CREATE TABLE IF NOT EXISTS faculty (
			id SERIAL PRIMARY KEY,
			school INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			is_deleted BOOLEAN NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := postgree.MainDBX.Exec(query)
	return err
}

type FacultySm struct {
	IdFaculty   int    `db:"id_faculty" json:"id_faculty"`
	NameFaculty string `db:"name_faculty" json:"name_faculty"`
}

type ListFacultySm []FacultySm

func (a *ListFacultySm) Marshall(c *gin.Context) []byte {
	res, err := json.Marshal(a)

	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка json.Marshal школ и предметов", "/directory")
		c.Abort()
		return []byte{}
	}
	return res
}

// получить список названий школ и факультетов по юзеру
func getListFacultySchoolByUser(user string) (ListFacultySm, error) {
	const query = `
	SELECT f.id AS id_faculty, f.name AS name_faculty
	FROM faculty AS f
	LEFT JOIN school AS s ON s.id = f.school
	WHERE s.director = $1
	AND f.is_deleted = false;
	`

	r := ListFacultySm{}
	err := postgree.MainDBX.Select(&r, query, user)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r, nil
}

func GetListFacultySmSchoolSmByUser(c *gin.Context, user string) (res ListFacultySm) {

	res, err := getListFacultySchoolByUser(user)

	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка списка школ и предметов", "/directory")
		c.Abort()
		return
	}
	return
}

func DeleteFaculty(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		webmessage.Err(c, fmt.Errorf("id not provided"), "ID факультета не указан", "/directory")
		c.Abort()
		return
	}

	query := `UPDATE faculty SET is_deleted = true, updated_at = $2 WHERE id = $1;`
	_, err := postgree.MainDB.Exec(query, id, time.Now())
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка при удалении факультета", "/directory")
		c.Abort()
		return
	}
}

func GetAllFacultysFromUser(c *gin.Context) (res ListFaculty) {
	const query = `
	SELECT f.id, s.name AS school_name, f.school, f.name, f.description,
		f.created_at, f.updated_at, f.is_deleted
	FROM faculty AS f
		LEFT JOIN school AS s ON s.id = f.school
	WHERE s.director = $1
	AND f.is_deleted = false;`

	session := sessions.Default(c)
	login := session.Get("user")

	err := postgree.MainDBX.Select(&res, query, login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка списка факультетов", "/directory")
		c.Abort()
		return
	}
	return
}

func GetFacultyByID(c *gin.Context) (res Faculty) {

	id := c.Query("id")
	if id == "" {
		webmessage.Err(c, fmt.Errorf("id not provided"), "ID факультета не указан", "/directory")
		c.Abort()
		return
	}

	const query = `
		SELECT id, school, name, description, created_at, updated_at, is_deleted
		FROM faculty
		WHERE id = $1;`

	err := postgree.MainDBX.Get(&res, query, id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка загрузки факультета из БД", "/directory")
		c.Abort()
		return
	}
	return
}

func (faculty *Faculty) Rec(c *gin.Context) {

	faculty.CreatedAt = time.Now()
	faculty.UpdatedAt = time.Now()

	const query = `
		INSERT INTO faculty (school, name, description, created_at, updated_at, is_deleted)
		VALUES (:school, :name, :description, :created_at, :updated_at, :is_deleted);
	`

	_, err := postgree.MainDBX.NamedExec(query, faculty)
	if err != nil {
		fmt.Println(err)
		if c == nil {
			fmt.Println("ошибка записи факультета")
		} else {
			webmessage.Err(c, nil, "ошибка записи факультета", "/list/faculty")
			c.Abort()
			return
		}
	}
}
