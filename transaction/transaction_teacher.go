package transaction

import (
	"fmt"
	"schoolonline/dict"
	"schoolonline/kassa"
	"schoolonline/postgree"
	"schoolonline/sendmail"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

// записываем родителя и обновляем поля ссылки
func ExecTransactionInsertUserAndTeacher(c *gin.Context, us dict.User) {
	err := transactionInsertUserAndTeacher(us)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка transaction teacher", "/directory")
		c.Abort()
		return
	}
	sendmail.SendEmailAccess(c, us.Email, us.Login, us.PassWord, "Учитель")
}

func transactionInsertUserAndTeacher(us dict.User) error {

	// Начинаем транзакцию
	tx, err := postgree.MainDBX.Beginx()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Откат транзакции в случае ошибки
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback()
		}
	}()

	err = us.HashPassword()
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = tx.NamedExec(dict.InsertUsersQuery, us)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// записываем родителя и обновляем поля ссылки
func ExecTransactionInsertUserAndTeacherAndUpdateLink(c *gin.Context, us dict.User, linkCode string) {
	err := transactionInsertUserAndTeacherAndUpdateLink(us, linkCode)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка transaction teacher", "/directory")
		c.Abort()
		return
	}
	sendmail.SendEmailAccess(c, us.Email, us.Login, us.PassWord, "Учитель")
}

func transactionInsertUserAndTeacherAndUpdateLink(us dict.User, linkCode string) error {

	// Начинаем транзакцию
	tx, err := postgree.MainDBX.Beginx()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Откат транзакции в случае ошибки
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback()
		}
	}()

	err = us.HashPassword()
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = tx.NamedExec(dict.InsertUsersQuery, us)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Обновляем поля active и complete в записи linkregistration
	const updateLinkQuery = `
		UPDATE linkregistration
		SET active = false, complete = true
		WHERE code = $1;
	`

	_, err = tx.Exec(updateLinkQuery, linkCode)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ExecTtransactionPaimentTeacher(c *gin.Context, login string, pay int) error {
	return transactionPaimentTeacher(login, pay)
}

func transactionPaimentTeacher(login string, pay int) error {

	// Начинаем транзакцию
	tx, err := postgree.MainDBX.Beginx()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Откат транзакции в случае ошибки
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback()
		}
	}()

	cr := time.Now()
	dl := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)

	// запись в кассу по teacher
	_, err = tx.Exec(kassa.QuweryInsertKassaByTeacherByTeacherLogin, login, pay,
		"3.2", "Выплата заработанных средств учителю", cr, dl, dl, false)
	if err != nil {
		//fmt.Println(4, err)
		return err
	}

	// списываем с баланса учителя по логину
	_, err = tx.Exec(dict.QueryUpdateUserBalanceMinusByLogin, pay, login)
	if err != nil {
		//fmt.Println(3, err)
		return err
	}

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
