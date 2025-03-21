package transaction

import (
	"fmt"
	"schoolonline/dict"
	"schoolonline/kassa"
	"schoolonline/postgree"
	"schoolonline/webmessage"
	"time"

	"github.com/gin-gonic/gin"
)

// записываем проведенный урок
func ExecTransactionLessonComplete(c *gin.Context, lesson dict.Lesson, user dict.User) {
	err := transactionLessonComplete(lesson, user)
	if err != nil {
		fmt.Println(0, err)
		webmessage.Err(c, err, "Ошибка transaction урок", "/directory")
		c.Abort()
		return
	}
}

// транзакции по проведенному уроку
func transactionLessonComplete(lesson dict.Lesson, user dict.User) error {

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

	lessonComplete := dict.LessonComplete{LessonId: lesson.ID, CreatorLogin: user.Login, CreateTime: time.Now(), IsDeleted: false}

	// записываем комплектность урока
	_, err = tx.NamedExec(dict.QueryLessonComplete, lessonComplete)
	if err != nil {
		//fmt.Println(1, err)
		return err
	}

	// снимаем с баланса родителя по логину студента
	_, err = tx.Exec(dict.QueryUpdateUserBalanceMinusBStudentLogin, lesson.CostLesson, lesson.StudentLogin)
	if err != nil {
		//fmt.Println(2, err)
		return err
	}

	// пишем на баланс учителя по логину
	_, err = tx.Exec(dict.QueryUpdateUserBalancePlusByLogin, lesson.TeacherBaseRate, lesson.TeacherLogin)
	if err != nil {
		//fmt.Println(3, err)
		return err
	}

	cr := time.Now()
	dl := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)

	// запись в кассу по родителю
	_, err = tx.Exec(kassa.QuweryInsertKassaByParentByStudentLogin, lesson.StudentLogin, lesson.CostLesson,
		"1.2", "Списание средств c родителя за проведенный урок", cr, dl, dl, false)
	if err != nil {
		//fmt.Println(4, err)
		return err
	}

	// запись в кассу по учителю
	_, err = tx.Exec(kassa.QuweryInsertKassaByTeacherByTeacherLogin, lesson.TeacherLogin, lesson.TeacherBaseRate,
		"3.1", "Зачисление средств учителю за проведенный урок", cr, dl, dl, false)
	if err != nil {
		//fmt.Println(5, err)
		return err
	}

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		//fmt.Println(6, err)
		return err
	}

	return nil
}
