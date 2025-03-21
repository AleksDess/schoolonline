package dict

import (
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ID          int       `db:"id" json:"id"`
	School      int       `db:"school" json:"school" form:"school"`
	Faculty     int       `db:"faculty" json:"faculty" form:"faculty"`
	Name        string    `db:"name" json:"name" form:"name"`
	Description string    `db:"description" json:"description" form:"description"`
	IsDeleted   bool      `db:"is_deleted" json:"is_deleted" form:"is_deleted"`
	CreatedAt   time.Time `db:"created_at" json:"created_at" form:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at" form:"updated_at"`
	FacultyName string    `db:"faculty_name"`
	SchoolName  string    `db:"school_name"`
}

type ListItem []Item

func CreateItemTableIfNotExists() error {
	const query = `
		CREATE TABLE IF NOT EXISTS item (
			id SERIAL PRIMARY KEY,
			school INTEGER NOT NULL,
			faculty INTEGER NOT NULL,
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

func GetAllItemsFromUser(c *gin.Context, user string) (res ListItem) {
	const query = `
	SELECT i.id, s.name AS school_name, f.name AS faculty_name, i.name,
		i.description, i.created_at, i.updated_at, i.is_deleted
	FROM item AS i
		LEFT JOIN school AS s ON s.id = i.school
		LEFT JOIN faculty AS f ON f.id = i.faculty
	WHERE s.director = $1
	AND i.is_deleted = false;
	`

	err := postgree.MainDBX.Select(&res, query, user)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка списка предметов", "/directory")
		c.Abort()
		return
	}

	return
}

// удалить предмет
func DeleteItem(c *gin.Context) {

	id := c.Query("id")
	if id == "" {
		webmessage.Err(c, fmt.Errorf("id not provided"), "ID факультета не указан", "/directory")
		c.Abort()
		return
	}

	query := `UPDATE item	SET is_deleted = true WHERE id = $1;`
	_, err := postgree.MainDB.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка при удалении предмета", "/directory")
		c.Abort()
		return
	}
}

func GetAllItems() (ListItem, error) {
	const query = `
		SELECT id, school, faculty, name, description, created_at, updated_at, is_deleted
		FROM item
		WHERE is_deleted = false;
	`

	var items ListItem
	err := postgree.MainDBX.Select(&items, query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return items, nil
}

func GetAllItemsBySchoolId(c *gin.Context, id string) (res ListItem) {
	const query = `
		SELECT id, school, faculty, name, description, created_at, updated_at, is_deleted
		FROM item
		WHERE school = $1
		AND is_deleted = false;`

	err := postgree.MainDBX.Select(&res, query, id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка чтения предметов школы", "/list/student")
		c.Abort()
		return
	}
	return
}

func GetItemByID(c *gin.Context) (res Item) {

	id := c.Query("id")
	if id == "" {
		webmessage.Err(c, fmt.Errorf("id not provided"), "ID предмета не указан", "/directory")
		c.Abort()
		return
	}

	const query = `
		SELECT id, school, faculty, name, description, created_at, updated_at, is_deleted
		FROM item
		WHERE id = $1;`

	err := postgree.MainDBX.Get(&res, query, id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, nil, "Ошибка загрузки предмета из БД", "/list/item")
		c.Abort()
		return
	}
	return
}

func (item *Item) Rec(c *gin.Context) {

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	const query = `
		INSERT INTO item (school, faculty, name, description, created_at, updated_at, is_deleted)
		VALUES (:school, :faculty, :name, :description, :created_at, :updated_at, :is_deleted);`

	_, err := postgree.MainDBX.NamedExec(query, item)
	if err != nil {
		fmt.Println(err)
		if c == nil {
			fmt.Println("ошибка записи предмета")
		} else {
			webmessage.Err(c, err, "Ошибка записи предмета", "/list/item")
			c.Abort()
		}

	}
}
