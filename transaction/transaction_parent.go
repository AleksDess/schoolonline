package transaction

import (
	"fmt"
	"schoolonline/dict"
	"schoolonline/postgree"
	"schoolonline/sendmail"
	"schoolonline/webmessage"

	"github.com/gin-gonic/gin"
)

// записываем учителя и обновляем поля ссылки
func ExecTransactionInsertUserAndParent(c *gin.Context, us dict.User) {
	err := transactionInsertUserAndParent(us)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка transaction parent", "/directory")
		c.Abort()
		return
	}
	sendmail.SendEmailAccess(c, us.Email, us.Login, us.PassWord, "Родитель")
}

func transactionInsertUserAndParent(us dict.User) error {

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

// записываем учителя и обновляем поля ссылки
func ExecTransactionInsertUserAndParentAndUpdateLink(c *gin.Context, us dict.User, linkCode string) {
	err := transactionInsertUserAndParentAndUpdateLink(us, linkCode)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "Ошибка transaction parent", "/directory")
		c.Abort()
		return
	}
	sendmail.SendEmailAccess(c, us.Email, us.Login, us.PassWord, "Родитель")
}

func transactionInsertUserAndParentAndUpdateLink(us dict.User, linkCode string) error {

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
