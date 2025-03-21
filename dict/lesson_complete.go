package dict

import (
	"schoolonline/postgree"
	"time"
)

// проведенное занятие
type LessonComplete struct {
	ID           int       `db:"id" json:"id" form:"id"`
	LessonId     int       `db:"lesson_id" json:"lesson_id" form:"lesson_id"`
	CreatorLogin string    `db:"creator_login" json:"creator_login" form:"creator_login"`
	CreateTime   time.Time `db:"create_time" json:"create_time" form:"create_time"`
	IsDeleted    bool      `db:"is_deleted" json:"is_deleted" form:"is_deleted"`
}

func CreateLessonCompleteTableIfNotExists() error {
	const query = `
        CREATE TABLE IF NOT EXISTS lesson_complete (
            id SERIAL PRIMARY KEY,
            lesson_id INT NOT NULL,
            creator_login TEXT NOT NULL,
			create_time TIMESTAMP NOT NULL,
            is_deleted BOOLEAN NOT NULL
        );
    `

	_, err := postgree.MainDBX.Exec(query)
	return err
}

var QueryLessonComplete = `INSERT INTO lesson_complete (
            lesson_id, creator_login, create_time, is_deleted
        )
        VALUES (
            :lesson_id, :creator_login, :create_time, :is_deleted
        );`
