package directory

import (
	"database/sql"
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

type LinkRegistration struct {
	Id          int       `db:"id"`
	SchoolId    int       `db:"school_id"`
	SchoolName  string    `db:"school_name"`
	Code        string    `db:"code"`
	Role        string    `db:"role"`
	Creator     string    `db:"creator"`
	CreatedTime time.Time `db:"created_time"`
	Active      bool      `db:"active"`
	Complete    bool      `db:"complete"`
}

type ListLinkRegistration []LinkRegistration

func CreateLinkRegistrationTableIfNotExists() error {
	const query = `
		CREATE TABLE IF NOT EXISTS linkregistration (
			id SERIAL PRIMARY KEY,
			school_id INTEGER NOT NULL,
			code TEXT NOT NULL,
			role TEXT NOT NULL,
			creator TEXT NOT NULL,
			created_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			active BOOLEAN DEFAULT TRUE,
			complete BOOLEAN DEFAULT FALSE
		);
	`

	_, err := postgree.MainDBX.Exec(query)
	return err
}

func GetLinkRegistrationByCode(code string) (*LinkRegistration, error) {
	const query = `
		SELECT l.id, l.school_id, l.code, l.role, l.creator, l.created_time, l.active, l.complete, s.name AS school_name
		FROM linkregistration AS l
		LEFT JOIN school AS s ON s.id = l.school_id
		WHERE code = $1;`

	var linkregistration LinkRegistration
	err := postgree.MainDBX.Get(&linkregistration, query, code)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, nil // Если директор не найден, возвращаем nil
		}
		return nil, err
	}

	return &linkregistration, nil
}

func (linkregistration *LinkRegistration) Rec(c *gin.Context) {

	linkregistration.CreatedTime = time.Now()
	linkregistration.Active = true
	linkregistration.Complete = false

	const query = `
		INSERT INTO linkregistration (school_id, code, role, creator, created_time, active, complete)
		VALUES (:school_id, :code, :role, :creator, :created_time, :active, :complete)
		RETURNING id;
	`

	stmt, err := postgree.MainDBX.PrepareNamed(query)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка записи linkregistration на регистрацию в БД", "/menu")
		c.Abort()
		return
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(linkregistration).Scan(&id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка записи linkregistration на регистрацию в БД", "/menu")
		c.Abort()
		return
	}

	linkregistration.Id = id
}
