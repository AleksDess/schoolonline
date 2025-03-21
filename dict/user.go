package dict

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// структура User
type User struct {
	Login         string        `db:"login" json:"login" form:"login"`
	Email         string        `db:"email" json:"email" form:"email"`
	Phone         string        `db:"phone" json:"phone" form:"phone"`
	PassWord      string        `db:"pass_word" json:"pass_word" form:"pass_word"`
	Creator       string        `db:"creator" json:"creator" form:"creator"`
	Hash          string        `db:"hash"`
	Role          string        `db:"role" json:"role" form:"role"`
	FirstName     string        `db:"first_name" json:"first_name"`
	LastName      string        `db:"last_name" json:"last_name"`
	SchoolID      int           `db:"school_id" json:"school_id" form:"school_id"`
	ParentLogin   string        `db:"parent_login" json:"parent_login" form:"parent_login"`
	CurrencyId    int           `db:"currency_id" json:"currency_id" form:"currency_id"`
	YearOfBirth   int           `db:"year_of_birth" json:"year_of_birth"`
	ListItem      pq.Int64Array `db:"list_item" json:"list_item" form:"list_item"`
	TgId          int64         `db:"tg_id" json:"tg_id" form:"tg_id"`
	VbId          string        `db:"vb_id" json:"vb_id" form:"vb_id"`
	Balance       int           `db:"balance" json:"balance" form:"balance"`
	PimentDetails string        `db:"piment_details" json:"piment_details" form:"piment_details"`
	UserPhoto     string        `db:"user_photo" json:"user_photo" form:"user_photo"`
	UserPhone     string        `db:"user_phone" json:"user_phone" form:"user_phone"`
	UserDevise    string        `db:"user_devise" json:"user_devise" form:"user_devise"`
	IsDeleted     bool          `db:"is_deleted" json:"is_deleted" form:"is_deleted"`
	VerifyEmail   bool          `db:"verify_email" json:"verify_email" form:"verify_email"`
	VerifyCode    string        `db:"verify_code" json:"verify_code" form:"verify_code"`
	CreateTime    time.Time     `db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime    time.Time     `db:"update_time" json:"update_time" form:"update_time"`
	VerifyTime    time.Time     `db:"verify_time" json:"verify_time" form:"verify_time"`
}

// функция создания таблицы
// постгресс User
func CreateTableUser() {
	_, err := postgree.MainDB.Exec(`
CREATE TABLE IF NOT EXISTS users (
  login TEXT PRIMARY KEY,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,
  pass_word TEXT NOT NULL,
  creator TEXT,
  hash TEXT NOT NULL,
  role TEXT NOT NULL,
  first_name TEXT,
  last_name TEXT,
  school_id INTEGER,
  parent_login TEXT,
  currency_id INTEGER,
  year_of_birth INTEGER,
  list_item INTEGER[],
  tg_id BIGINT NOT NULL,
  vb_id TEXT NOT NULL,
  balance INTEGER NOT NULL,
  piment_details TEXT,
  user_photo TEXT,
  user_phone TEXT,
  user_devise TEXT,
  is_deleted BOOLEAN DEFAULT FALSE,
  verify_email BOOLEAN DEFAULT FALSE,
  verify_code TEXT,
  create_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  update_time TIMESTAMPTZ,
  verify_time TIMESTAMPTZ
);
`)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error creating users table: %v", err)
	}
}

// список элементов User
type ListUser []User

const InsertUsersQuery = `
	INSERT INTO users
	(login, email, phone, pass_word, creator, hash, role, first_name, last_name, school_id, parent_login, currency_id, year_of_birth, list_item, 
	tg_id, vb_id, balance, piment_details, user_photo, user_phone, user_devise, is_deleted, verify_email, verify_code, 
	create_time, update_time, verify_time)
	VALUES 
	(:login, :email, :phone, :pass_word, :creator, :hash, :role, :first_name, :last_name, :school_id, :parent_login, :currency_id, :year_of_birth, :list_item, 
	:tg_id, :vb_id, :balance, :piment_details, :user_photo, :user_phone, :user_devise, :is_deleted, :verify_email, :verify_code, 
	:create_time, :update_time, :verify_time)
`

func (a *User) Rec(c *gin.Context) {
	var err error

	// Генерация хэша пароля
	a.Hash, err = hashPassword(a.PassWord)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка записи юзера", "/directory")
		c.Abort()
		return
	}

	// Вставка записи в базу данных
	_, err = postgree.MainDBX.NamedExec(InsertUsersQuery, a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка записи юзера", "/directory")
		c.Abort()
		return
	}
}

const ReadUserQuery = `SELECT login, email, phone, pass_word, creator, hash, role, first_name, last_name, school_id, parent_login, 
	       currency_id, year_of_birth, list_item, tg_id, vb_id, balance, piment_details, user_photo, user_phone, user_devise, 
	       is_deleted, verify_email, verify_code, create_time, update_time, verify_time
	FROM users `

// Функция для сканирования строки в структуру User
func scanUser(row *sql.Row, a *User) error {
	return row.Scan(
		&a.Login, &a.Email, &a.Phone, &a.PassWord, &a.Creator, &a.Hash, &a.Role, &a.FirstName, &a.LastName,
		&a.SchoolID, &a.ParentLogin, &a.CurrencyId, &a.YearOfBirth, &a.ListItem, &a.TgId,
		&a.VbId, &a.Balance, &a.PimentDetails, &a.UserPhoto, &a.UserPhone, &a.UserDevise, &a.IsDeleted, &a.VerifyEmail,
		&a.VerifyCode, &a.CreateTime, &a.UpdateTime, &a.VerifyTime,
	)
}

func GetUserByVerifyCode(verifyCode string) (a User, err error) {
	query := ReadUserQuery + `
	WHERE verify_code = $1 
	LIMIT 1`
	row := postgree.MainDBX.QueryRow(query, verifyCode)
	err = scanUser(row, &a)
	return
}

// обновление после верификации
func UpdateUserVerifyStatus(verifyCode string) (error, bool) {
	// SQL-запрос для обновления полей VerifyEmail и VerifyTime
	query := `UPDATE users 
              SET verify_email = TRUE, verify_time = NOW(), update_time = NOW()
              WHERE verify_code = $1 AND verify_email = FALSE`

	// Выполнение запроса
	result, err := postgree.MainDBX.Exec(query, verifyCode)
	if err != nil {
		fmt.Println(err)
		return err, false
	}

	// Проверяем, было ли обновлено хотя бы одно поле
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err, false
	}

	if rowsAffected == 0 {
		return nil, false
	}

	return nil, true
}

// чтение элемента User
func GetUserByLogin(login string) (a User, err error) {
	query := `SELECT * FROM users WHERE login = $1 LIMIT 1`
	row := postgree.MainDBX.QueryRow(query, login)
	err = scanUser(row, &a)
	return
}

// чтение элемента User по tg_id
func GetUserByTgId(id int64) (a User, err error) {
	query := `SELECT * FROM users WHERE tg_id = $1 LIMIT 1`
	row := postgree.MainDBX.QueryRow(query, id)
	err = scanUser(row, &a)
	return
}

// чтение элемента User по vb_id
func GetUserByVbId(id int64) (a User, err error) {
	query := `SELECT * FROM users WHERE vb_id = $1 LIMIT 1`
	row := postgree.MainDBX.QueryRow(query, id)
	err = scanUser(row, &a)
	return
}

// Проверка наличия пользователя по login и получение TgId
func UserExists(login string) (bool, int64, error) {
	var tgId int64
	query := `SELECT tg_id FROM users WHERE login = $1 LIMIT 1`
	err := postgree.MainDBX.QueryRow(query, login).Scan(&tgId)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return false, 0, nil // Запись не найдена
		}
		return false, 0, err // Произошла ошибка
	}
	return true, tgId, nil // Запись найдена, возвращаем значение tg_id
}

// Обновление school_id пользователя
func updateUserSchoolId(tx *sqlx.Tx, login string, new int) error {
	query := `UPDATE users SET school_id = $1, update_time = $3 WHERE login = $2`
	_, err := tx.Exec(query, new, login, time.Now())
	return err
}

type UserBalance struct {
	Name    string `db:"name"`
	Balance int    `db:"balance"`
	Code    string `db:"code"`
}

// получить баланс по vb ид
func GetUserBalanceByTgId(id int64) (res UserBalance, err error) {
	queryCheck := `
		SELECT 
			p.first_name || ' ' || p.last_name AS name, 
			p.balance, 
			c.code 
		FROM users AS p
		LEFT JOIN currency AS c ON c.id = p.currency_id
		WHERE p.tg_id = $1`

	err = postgree.MainDBX.Get(&res, queryCheck, id)
	return
}

// получить баланс по vb ид
func GetUserBalanceByVbId(id string) (res UserBalance, err error) {
	queryCheck := `
		SELECT 
			p.first_name || ' ' || p.last_name AS name, 
			p.balance, 
			c.code 
		FROM users AS p
		LEFT JOIN currency AS c ON c.id = p.currency_id
		WHERE p.vb_id = $1`

	err = postgree.MainDBX.Get(&res, queryCheck, id)
	return
}

// обновить тг ид
func UpdateUserTgId(login string, newTgId int64) error {
	// Проверяем, существует ли уже запись с таким TgId и is_deleted = false
	var existingLogin string
	queryCheck := `SELECT login FROM users WHERE tg_id = $1 AND is_deleted = false`
	err := postgree.MainDBX.Get(&existingLogin, queryCheck, newTgId)
	if err == nil {
		return fmt.Errorf("TgId %d уже используется для пользователя с логином %s", newTgId, existingLogin)
	} else if err != sql.ErrNoRows {
		// Если ошибка не связана с отсутствием записи, возвращаем ее
		return fmt.Errorf("ошибка проверки существующего TgId: %v", err)
	}

	// Если проверка прошла, обновляем TgId
	queryUpdate := `UPDATE users SET tg_id = $1, update_time = $3 WHERE login = $2`
	_, err = postgree.MainDBX.Exec(queryUpdate, newTgId, login, time.Now())
	if err != nil {
		return fmt.Errorf("ошибка обновления TgId для логина %s: %v", login, err)
	}

	return nil
}

// обновить вайбер ид
func UpdateUserVbId(login string, newVbId string) error {
	// Проверяем, существует ли уже запись с таким VbId и is_deleted = false
	var existingLogin string
	queryCheck := `SELECT login FROM users WHERE vb_id = $1 AND is_deleted = false`
	err := postgree.MainDBX.Get(&existingLogin, queryCheck, newVbId)
	if err == nil {
		return fmt.Errorf("VbId %s уже используется для пользователя с логином %s", newVbId, existingLogin)
	} else if err != sql.ErrNoRows {
		// Если ошибка не связана с отсутствием записи, возвращаем ее
		return fmt.Errorf("ошибка проверки существующего VbId: %v", err)
	}

	// Если проверка прошла, обновляем VbId
	queryUpdate := `UPDATE users SET vb_id = $1, update_time = $3 WHERE login = $2`
	_, err = postgree.MainDBX.Exec(queryUpdate, newVbId, login, time.Now())
	if err != nil {
		return fmt.Errorf("ошибка обновления VbId для логина %s: %v", login, err)
	}

	return nil
}

// Обновление Email пользователя
func UpdateUserEmail(login string, newEmail string) error {
	query := `UPDATE users SET email = $1, update_time = $3 WHERE login = $2`
	_, err := postgree.MainDBX.Exec(query, newEmail, login, time.Now())
	return err
}

const QueryUpdateUserBalancePlusByTgId = `UPDATE users SET balance = balance + $1 WHERE tg_id = $2`
const QueryUpdateUserBalanceMinusByTgId = `UPDATE users SET balance = balance - $1 WHERE tg_id = $2`

const QueryUpdateUserBalancePlusByLogin = `UPDATE users SET balance = balance + $1 WHERE login = $2`
const QueryUpdateUserBalanceMinusByLogin = `UPDATE users SET balance = balance - $1 WHERE login = $2`

const QueryUpdateUserBalanceMinusBStudentLogin = `
UPDATE users 
SET balance = balance - $1 
WHERE login = (SELECT parent_login FROM users WHERE login = $2);
`

// Обновление Hash пользователя
func UpdateUserHash(login string, newHash string) error {
	query := `UPDATE users SET hash = $1, update_time = $3 WHERE login = $2`
	_, err := postgree.MainDBX.Exec(query, newHash, login, time.Now())
	return err
}

// чтение списка элементов User
func (a *ListUser) GetAll() error {
	err := postgree.MainDBX.Select(a, "SELECT * FROM user")
	return err
}

// чтение списка элементов User
func GetAllListUser() (a ListUser, err error) {
	err = postgree.MainDBX.Select(&a, "SELECT * FROM users")
	return
}

// UpdateUserPassWord обновляет PassWord заданным Login.
func (a *User) UpdateUserPassWord(new string) error {

	q := "UPDATE users SET password = $1, update_time = $3 WHERE login = $2"
	_, err := postgree.MainDB.Query(q, new, a.Login, time.Now())

	return err
}

func (a *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.PassWord), 14)
	if err != nil {
		fmt.Println(err)
		return err
	}
	a.Hash = string(hashedPassword)
	return nil
}

// Функция для хеширования пароля
func hashPassword(password string) (string, error) {
	// Генерируем хеш пароля с указанной "стоимостью" (рекомендуется использовать значение 14)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(hashedPassword), nil
}

func (a *ListUser) Marshall(c *gin.Context) string {
	r, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка json.Marshal школ и предметов", "/directory")
		c.Abort()
		return ""
	}
	return string(r)
}
