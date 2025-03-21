package dict

import (
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

type School struct {
	ID          int       `db:"id" json:"id" form:"id"`
	Name        string    `db:"name" json:"name" form:"name"`
	Description string    `db:"description" json:"description" form:"description"`
	Website     string    `db:"website" json:"website" form:"website"`
	Email       string    `db:"email" json:"email" form:"email"`
	Phone       string    `db:"phone" json:"phone" form:"phone"`
	Address     string    `db:"address" json:"address" form:"address"`
	Director    string    `db:"director" json:"director" form:"director"`
	IsDeleted   bool      `db:"is_deleted" json:"is_deleted" form:"is_deleted"`
	CreatedAt   time.Time `db:"created_at" json:"created_at" form:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at" form:"updated_at"`
}

type ListSchool []School

func CreateSchoolTableIfNotExists() error {
	const query = `
		CREATE TABLE IF NOT EXISTS school (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			website TEXT,
			email TEXT,
			phone TEXT,
			address TEXT,
			director TEXT,
			is_deleted BOOLEAN NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := postgree.MainDBX.Exec(query)
	return err
}

// список школ
func GetAllSchoolsFromUser(c *gin.Context, login string) (res ListSchool) {
	const query = `
		SELECT id, name, description, website, email, phone, address,
		director, is_deleted, created_at, updated_at FROM school
		WHERE director = $1
		AND is_deleted = false;`

	err := postgree.MainDBX.Select(&res, query, login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка списка школ", "/directory")
		c.Abort()
		return
	}

	return
}

// удалить школу
func DeleteSchool(c *gin.Context, id string) {
	query := `UPDATE school	SET is_deleted = true WHERE id = $1;`
	_, err := postgree.MainDB.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка при удалении школы", "/directory")
		c.Abort()
		return
	}
}

// получить одну школу
func GetSchoolByID(c *gin.Context, id string) (res School) {
	const query = `
		SELECT id, name, description, website, email, phone, address, director, is_deleted, created_at, updated_at
		FROM school
		WHERE id = $1;`

	err := postgree.MainDBX.Get(&res, query, id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка загрузки школы из БД", "/directory")
		c.Abort()
		return
	}
	return
}

// получить одну школу  по юзеру
func GetSchoolByUser(c *gin.Context, user string) School {
	const query = `
		SELECT id, name, description, website, email, phone, address, director, is_deleted, created_at, updated_at
		FROM school
		WHERE director = $1
		AND is_deleted = false
		LIMIT 1;
	`

	var school School
	err := postgree.MainDBX.Get(&school, query, user)
	if err != nil {
		fmt.Println(err)
		webmessage.SendMessage(c, "У Вас не добавлена ни одна школа", "/input/school")
		c.Abort()
		return school
	}

	return school
}

// Структура для хранения результата запроса
type SchoolSmall struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// Получить ID школы по имени пользователя
func GetIdSchoolByUser(c *gin.Context, user string) (res SchoolSmall) {
	const query = `
		SELECT id, name FROM school
		WHERE director = $1
		AND is_deleted = false
		LIMIT 1;	`

	err := postgree.MainDBX.Get(&res, query, user)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Не удалось получить ID школы", "/directory")
		c.Abort()
		return
	}

	return
}
func (school *School) Rec(c *gin.Context, director string) {
	err := school.rec(director)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка списка школы", "/directory")
		c.Abort()
		return
	}
}

const insertSchoolQuery = `
		INSERT INTO school (name, description, website, email, phone, address, director, is_deleted, created_at, updated_at)
		VALUES (:name, :description, :website, :email, :phone, :address, :director, :is_deleted, :created_at, :updated_at)
		RETURNING id;
	`

func (school *School) rec(director string) error {
	school.CreatedAt = time.Now()
	school.UpdatedAt = time.Now()
	school.Director = director

	// Начинаем транзакцию
	tx, err := postgree.MainDBX.Beginx()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println(err)
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	stmt, err := tx.PrepareNamed(insertSchoolQuery)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(school).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	school.ID = id

	// Обновляем ID школы у пользователя
	err = updateUserSchoolId(tx, director, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
