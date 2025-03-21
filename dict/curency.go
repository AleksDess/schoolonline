package dict

import (
	"fmt"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"sort"

	"github.com/gin-gonic/gin"
)

type Currency struct {
	ID        int    `db:"id" json:"id"`
	SchoolId  int    `db:"school_id"`
	Code      string `db:"code" json:"code"`
	Symbol    string `db:"symbol" json:"symbol"`
	Name      string `db:"name" json:"name"`
	RusName   string `db:"rus_name" json:"rus_name"`
	IsDeleted bool   `db:"is_deleted" json:"is_deleted" form:"is_deleted"`
}

type ListCurrency []Currency

var listCurrency = map[int]Currency{
	1: {1, 0, "USD", "$", "United States Dollar", "Доллар США", false},
	2: {2, 0, "EUR", "€", "Euro", "Евро", false},
	3: {3, 0, "UAH", "₴", "Ukrainian Hryvnia", "Украинская гривна", false},
	4: {4, 0, "RUB", "₽", "Russian Ruble", "Российский рубль", false},
	5: {5, 0, "BYN", "Br", "Belarusian Ruble", "Белорусский рубль", false},
	6: {6, 0, "GBP", "£", "British Pound Sterling", "Британский фунт стерлингов", false},
	7: {7, 0, "CHF", "CHF", "Swiss Franc", "Швейцарский франк", false},
	8: {8, 0, "CNY", "¥", "Chinese Yuan", "Китайский юань", false},
}

func tableExists(tableName string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_name = $1
		);
	`

	var exists bool
	err := postgree.MainDB.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return exists, nil
}

func createCurrencyTable() error {
	const query = `
		CREATE TABLE currency (
			id SERIAL PRIMARY KEY,
			school_id INTEGER NOT NULL,
			code TEXT NOT NULL,
			symbol TEXT NOT NULL,
			name TEXT NOT NULL,
			rus_name TEXT NOT NULL,
			is_deleted BOOLEAN NOT NULL
		);`

	_, err := postgree.MainDB.Exec(query)
	return err
}

func (a *Currency) Rec(c *gin.Context) {
	const query = `
		INSERT INTO currency (school_id, code, symbol, name, rus_name, is_deleted)
		VALUES (:school_id, :code, :symbol, :name, :rus_name, :is_deleted);
	`
	_, err := postgree.MainDBX.NamedExec(query, a)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка добавления новой валюты", "/list/currency")
		c.Abort()
		return
	}
}

func CreateTablecurrency() {

	// Проверка существования таблицы currency
	exists, err := tableExists("currency")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to check if table exists:", err)
		return
	}

	if !exists {
		// Создание таблицы currency, если она не существует
		err = createCurrencyTable()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to create currency table:", err)
			return
		}

		// Вставка всех валют в таблицу
		err = insertCurrency()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to insert currency:", err)
			return
		}

		fmt.Println("Currency inserted successfully")
	}
}

func insertCurrency() error {

	rs := ListCurrency{}

	for _, i := range listCurrency {
		rs = append(rs, i)
	}

	sort.Slice(rs, func(i, j int) bool {
		return rs[i].ID < rs[j].ID
	})

	var c *gin.Context

	for _, i := range rs {
		i.Rec(c)
	}

	return nil
}

func GetAllCurrencyBySchoolId(c *gin.Context, id int) ListCurrency {
	const query = `
		SELECT id, school_id, code, symbol, name, rus_name, is_deleted
		FROM currency
		WHERE (school_id = $1 OR school_id = 0)
		AND is_deleted = false
		ORDER BY id;
	`

	var currency []Currency
	err := postgree.MainDBX.Select(&currency, query, id)
	if err != nil {
		fmt.Println(err)
		if c == nil {
			fmt.Println(err, "Ошибка чтения списка валют")
		} else {
			webmessage.Err(c, err, "Ошибка чтения списка валют", "/directory")
			c.Abort()
			return nil
		}
	}

	return currency
}

func GetAllCurrencyByUser(c *gin.Context, user User) (res ListCurrency) {

	const query = `
		SELECT c.id, c.school_id, c.code, c.symbol, c.name, c.rus_name, c.is_deleted
		FROM currency AS c
			LEFT JOIN school AS s ON s.id = c.school_id 
			WHERE (s.director = $1 OR c.school_id = 0) 
			AND c.is_deleted = false
		ORDER BY id
	`

	err := postgree.MainDBX.Select(&res, query, user.Login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка чтения БД", "/menu")
		c.Abort()
		return nil
	}

	return

}

func GetCurrencyByID(c *gin.Context) (res Currency) {

	id := c.Query("id")
	if id == "" {
		webmessage.Err(c, fmt.Errorf("id not provided"), "ID currency не указан", "/directory")
		c.Abort()
		return
	}

	const query = `
		SELECT id, school_id, code, symbol, name, rus_name, is_deleted
		FROM currency
		WHERE id = $1;	`

	err := postgree.MainDBX.Get(&res, query, id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ошибка загрузки currency из БД", "/directory")
		c.Abort()
		return
	}

	return
}

func GetCurrencyByUserLogin(login string) (res Currency, err error) {

	const query = `
		SELECT c.id, c.school_id, c.code, c.symbol, c.name, c.rus_name, c.is_deleted
		FROM currency AS c
		LEFT JOIN users AS u ON c.id = u.currency_id
		WHERE u.login = $1;	`

	err = postgree.MainDBX.Get(&res, query, login)
	return
}

// удалить student
func DeleteCurrency(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		webmessage.Err(c, fmt.Errorf("id not provided"), "ID currency не указан", "/directory")
		c.Abort()
		return
	}

	query := `UPDATE currency SET is_deleted = true WHERE id = $1;`
	_, err := postgree.MainDB.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка при удалении currencyt", "/directory")
		c.Abort()
		return
	}
}
