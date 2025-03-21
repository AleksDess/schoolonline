package postgree

import (
	"database/sql"
	"fmt"
	"schoolonline/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var MainDB *sql.DB
var MainDBX *sqlx.DB

func Set_idle_in_transaction_session_timeout() {
	q := `SET idle_in_transaction_session_timeout TO '300000'`
	_, err := MainDB.Exec(q)
	if err != nil {

		fmt.Println(err)
		return
	}
}

func RunDB() (err error) {
	MainDB, err = GetDB()
	if err != nil {

		fmt.Println(err)
		return
	}
	MainDBX, err = GetDBX()
	if err != nil {

		fmt.Println(err)
		return
	}
	return
}

func GetDB() (db *sql.DB, err error) {

	connStr := config.C.ConnServ

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		fmt.Println("err postgree not connect", err)
		return
	}
	return
}

func GetDBX() (db *sqlx.DB, err error) {

	connStr := config.C.ConnServ

	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		fmt.Println("err postgree not connect", err)
		return
	}
	return
}
