package createtable

import (
	"schoolonline/dict"
	"schoolonline/kassa"
	"schoolonline/metrics"
	"schoolonline/routes/directory"
	"schoolonline/tgbot/step"
)

func GoToCreateTable() {}

func CreateTable() {

	// user
	dict.CreateTableUser()
	// dict
	dict.CreateTablecurrency()
	dict.CreateSchoolTableIfNotExists()
	dict.CreateFacultyTableIfNotExists()
	dict.CreateItemTableIfNotExists()
	dict.CreateLessonTableIfNotExists()
	dict.CreateLessonCompleteTableIfNotExists()
	kassa.CreateTableKassa()
	// directory
	directory.CreateLinkRegistrationTableIfNotExists()
	// metrics
	metrics.CreateMetricsTableIfNotExists()
	// bot step
	step.CreateTableStepIsNotExist()

}

func CreateManager() {

	// man := manager.Manager{}
	// man.Id = "director"
	// man.FullName = "Богдан Ямщиков"
	// man.Login = "director"
	// man.Password = "ptihw1y8AB4PzjI"
	// man.TgId = 1001083516

	// man.RecNoConflict()

	// man.Id = "it"
	// man.FullName = "Александр Денисенко"
	// man.Login = "20dess20@gmail.com"
	// man.Password = "s3MatdDu5jiovw5"
	// man.TgId = 1791461936

	// man.RecNoConflict()

}
