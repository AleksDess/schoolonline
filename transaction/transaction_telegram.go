package transaction

import (
	"fmt"
	"schoolonline/dict"
	"schoolonline/kassa"
	"schoolonline/postgree"
	"time"
)

// Зачисление средств от родителя через телеграмм бот
func TransactionUserKassa(confirmedTgId, userTgId int64, cost int) error {

	// Начинаем транзакцию
	tx, err := postgree.MainDBX.Beginx()
	if err != nil {
		fmt.Println(0, err)
		return err
	}

	// Откат транзакции в случае ошибки
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback()
		}
	}()

	// пополняем баланс папе
	_, err = tx.Exec(dict.QueryUpdateUserBalancePlusByTgId, cost, userTgId)
	if err != nil {
		fmt.Println(1, err)
		return err
	}

	// формируем запись в кассу
	cr := time.Now()
	dl := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
	comment := "Зачисление средств от родителя через телеграмм бот"
	_, err = tx.Exec(kassa.QuweryInsertKassaByTelegram,
		userTgId, confirmedTgId, cost, "1.1", comment, cr, dl, dl, false)
	if err != nil {
		fmt.Println(2, err)
		return err
	}

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		fmt.Println(3, err)
		return err
	}

	return nil
}
