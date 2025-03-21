package step

import (
	"database/sql"
	"errors"
	"fmt"
	"schoolonline/postgree"
	"time"
)

// список функций
// ar  --  пополнение счета родителем

type Step struct {
	IdUser     int64     `db:"id_user"`
	Function   string    `db:"function"`
	Step       int       `db:"step"`
	Data       []byte    `db:"data"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

// функция создания таблицы
// постгресс Step
func CreateTableStepIsNotExist() {
	_, err := postgree.MainDBX.Exec(`
	CREATE TABLE IF NOT EXISTS bot_step (
		id_user  BIGSERIAL PRIMARY KEY,
		function VARCHAR(10),
		step INT,
		data   JSONB,
		create_time TIMESTAMP,
		update_time TIMESTAMP
		);`)

	if err != nil {

		fmt.Println(err)
	}
}

// печать элемента Step
func (a *Step) Print() {
	fmt.Println("")
	fmt.Println("----------------------------")
	fmt.Println("структура: Step")
	fmt.Println("Parent:      ", a.IdUser)
	fmt.Println("Step:        ", a.Step)
	fmt.Println("Function:    ", a.Function)
	fmt.Println("Data:        ", string(a.Data))
	fmt.Println("CreateTime:  ", a.CreateTime)
	fmt.Println("UpdateTime:  ", a.UpdateTime)
}

func GetStep(id int64) (res Step, found bool, err error) {
	query := `SELECT * FROM bot_step WHERE id_user = $1`
	err = postgree.MainDBX.Get(&res, query, id)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return res, false, nil
		}
		return res, false, err
	}
	return res, true, nil
}

// запись элемента   Step
func (a *Step) Rec() (err error) {

	a.CreateTime = time.Now()
	a.UpdateTime = time.Now()

	const query = `
INSERT INTO bot_step
(id_user, function, step, data, create_time, update_time)
VALUES 
(:id_user, :function, :step, :data, :create_time, :update_time)`

	a.CreateTime = time.Now()

	_, err = postgree.MainDBX.NamedExec(query, a)
	return err

}

func Get(id int64) (step Step, err error) {
	query := `SELECT * FROM bot_step WHERE id_user = $1`
	row := postgree.MainDBX.QueryRow(query, id)
	err = row.Scan(&step.IdUser, &step.Step, &step.Data, &step.CreateTime, &step.UpdateTime)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

// update
func (a *Step) Update() error {
	q := "UPDATE bot_step SET step = $1, data = $2, update_time = $4 WHERE id_user = $3"
	_, err := postgree.MainDBX.Query(q, a.Step, a.Data, a.IdUser, time.Now())
	fmt.Println(err)
	return err
}

// / delete
func Delete(id int64) error {
	q := "DELETE FROM bot_step WHERE id_user = $1"
	_, err := postgree.MainDBX.Exec(q, id)
	return err
}

func (a *Step) Delete(id int64) error {
	q := "DELETE FROM bot_step WHERE id_user = $1"
	_, err := postgree.MainDBX.Exec(q, id)
	return err
}

func DeleteAll() error {
	q := "DELETE FROM bot_step"
	_, err := postgree.MainDBX.Exec(q)
	return err
}
