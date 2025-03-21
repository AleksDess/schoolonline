package routes

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"schoolonline/config"
	"schoolonline/routes/access"
	"schoolonline/routes/blockip"
	"schoolonline/routes/directory"
	"schoolonline/routes/menu"
	"schoolonline/routes/message"
	"schoolonline/routes/registration"
	"schoolonline/routes/registration/entrance"
	"schoolonline/routes/startpage"
	"schoolonline/routes/user"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func GoToRunRoutes() {}

func RunRoutes() {

	gob.Register(registration.RegistrationData{})

	r := gin.New()

	r.SetFuncMap(template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("02.01.2006 15:04")
		},
		"formatDateSm": func(t time.Time) string {
			return t.Format("02.01.2006")
		},
		"formatTime": func(t time.Time) string {
			return t.Format("15:04:05")
		},
		"formatFloat64": func(t float64) string {
			return fmt.Sprintf("%.2f", t)
		},
		"toUpperCase": func(s string) string {
			return strings.ToUpper(s)
		},
		"toLowerCase": func(s string) string {
			return strings.ToLower(s)
		},
	})

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if config.C.Launch == "server" {
		gin.SetMode(gin.ReleaseMode)
	}

	err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Error setting trusted proxies: %v\n", err)
	}

	r.LoadHTMLFiles(pathHTML...)

	r.Static("/static", "./templates")

	secret := config.C.SessionKey

	store := cookie.NewStore([]byte(secret))

	store.Options(sessions.Options{
		MaxAge:   3600 * 24 * 7, // Время жизни сессии 7 day
		Path:     "/",           // Путь для сессий
		HttpOnly: true,          // Ограничиваем доступ только HTTP
		// Secure:   true,      // Включить только по HTTPS
	})

	r.Use(sessions.Sessions("mysession", store))

	r.Use(blockip.RequestTrackerMiddleware())

	r.GET("/", AuthRequired(), startPageHandler)
	r.GET("/index", indexHandler)
	r.GET("/help", AuthRequired(), helpHandler)
	r.GET("/about", AuthRequired(), aboutHandler)

	//  user
	r.GET("/user/profile", user.GetUserProfileHandler)
	r.GET("/user/setting", user.GetUserSettingHandler)
	r.POST("/user/profile", user.PostUserProfileHandler)
	r.POST("/user/setting", user.PostUserSettingHandler)
	r.POST("/user/edit/email", user.PostUserEditEmailHandler)
	r.POST("/user/edit/tg", user.PostUserEditTgHandler)

	// registration
	r.GET("/register", registration.RegisterGetHandler)
	r.POST("/register", registration.RegisterPostHandler)
	r.GET("/capcha", registration.CapchaGetHandler)
	r.POST("/capcha", registration.CapchaPostHandler)
	r.GET("/entrance", entrance.EntranceGetHandler)
	r.POST("/entrance", entrance.EntrancePostHandler)
	r.GET("/checkemail/:param", registration.ConfirmationOfRegisterPostHandler)

	// start page
	r.GET("/start/page/parent", AuthRequired(), startpage.GetStartPageHandler)
	r.POST("/start/page/parent", AuthRequired(), startpage.PostStartPageHandler)

	// message
	r.GET("/message/director", AuthRequired(), message.GetMessageDirectorHandler)

	//menu
	r.GET("/menu", AuthRequired(), menu.MenuGetHandler)
	r.GET("/logout", AuthRequired(), entrance.GetLogoutHandler)
	r.POST("/logout", AuthRequired(), entrance.PostLogoutHandler)
	r.GET("/check/username", registration.CheckUsernameHandler)
	r.GET("/directory", AuthRequired(), directory.GetDirectoryHandler)
	r.GET("/menu/links", AuthRequired(), menu.GetMenuLinksHandler)
	r.GET("/directory/links", AuthRequired(), directory.GetDirectoryLinksHandler)

	// school
	r.GET("/list/school", AuthRequired(), directory.GetDirectoryListSchoolHandler)
	r.GET("/input/school", AuthRequired(), directory.GetDirectoryInputSchoolHandler)
	r.POST("/input/school", AuthRequired(), directory.PostDirectoryInputSchoolHandler)
	r.GET("/view/school", AuthRequired(), directory.GetDirectoryViewSchoolHandler)
	r.DELETE("/delete/school", AuthRequired(), directory.GetDirectoryDeleteSchoolHandler)

	// pay
	r.GET("/pay", AuthRequired(), directory.GetDirectoryListPayHandler)

	//faculty
	r.GET("/list/faculty", AuthRequired(), directory.GetDirectoryListFacultyHandler)
	r.GET("/input/faculty", AuthRequired(), directory.GetDirectoryInputFacultyHandler)
	r.POST("/input/faculty", AuthRequired(), directory.PostDirectoryInputFacultyHandler)
	r.GET("/view/faculty", AuthRequired(), directory.GetDirectoryViewFacultyHandler)
	r.DELETE("/delete/faculty", AuthRequired(), directory.GetDirectoryDeleteFacultyHandler)

	// item
	r.GET("/list/item", AuthRequired(), directory.GetDirectoryListItemHandler)
	r.GET("/input/item", AuthRequired(), directory.GetDirectoryInputItemHandler)
	r.POST("/input/item", AuthRequired(), directory.PostDirectoryInputItemHandler)
	r.GET("/view/item", AuthRequired(), directory.GetDirectoryViewItemHandler)
	r.DELETE("/delete/item", AuthRequired(), directory.GetDirectoryDeleteItemHandler)

	// parent
	r.GET("/list/parent", AuthRequired(), directory.GetDirectoryListParentHandler)
	r.GET("/input/parent", AuthRequired(), directory.GetDirectoryInputParentHandler)
	r.POST("/input/parent", directory.PostDirectoryInputParentHandler)
	r.GET("/input/parent/create/link/:param", AuthRequired(), directory.GetDirectoryInputParentLinkHandler)
	r.GET("/input/parent/registration/:param", directory.GetDirectoryInputParentRegistration)
	r.GET("/view/parent", AuthRequired(), directory.GetDirectoryViewParentHandler)
	r.GET("/refill/parent/balance", AuthRequired(), directory.GetRefillParentBalanceHandler)
	r.DELETE("/delete/parent", AuthRequired(), directory.GetDirectoryDeleteParentHandler)

	// student
	r.GET("/list/student", AuthRequired(), directory.GetDirectoryListStudentHandler)
	r.POST("/input/student", AuthRequired(), directory.PostDirectoryInputStudentHandler)
	r.GET("/view/student", AuthRequired(), directory.GetDirectoryViewStudentHandler)
	r.DELETE("/delete/student", AuthRequired(), directory.GetDirectoryDeleteStudentHandler)

	// lesson
	r.POST("/input/lesson", AuthRequired(), directory.PostDirectoryInputLessonHandler)
	r.GET("/input/lesson", AuthRequired(), directory.GetDirectoryInputLessonHandler)
	r.POST("/save/lesson", AuthRequired(), directory.PostDirectorySaveLessonHandler)
	r.GET("/list/lesson", AuthRequired(), directory.GetListLessonHandler)
	r.GET("/view/lesson", AuthRequired(), directory.GetViewLessonHandler)
	r.GET("/lesson/cancel", AuthRequired(), directory.GetLessonCancelHandler)
	r.GET("/lesson/complete", AuthRequired(), directory.GetLessonCompleteHandler)

	// teacher
	r.GET("/list/teacher", AuthRequired(), directory.GetDirectoryListTeacherHandler)
	r.GET("/input/teacher", AuthRequired(), directory.GetDirectoryInputTeacherHandler)
	r.POST("/input/teacher", directory.PostDirectoryInputTeacherHandler)
	r.GET("/input/teacher/create/link/:param", AuthRequired(), directory.GetDirectoryInputTeacherLinkHandler)
	r.GET("/input/teacher/registration/:param", directory.GetDirectoryInputTeacherRegistration)
	r.GET("/view/teacher", AuthRequired(), directory.GetDirectoryViewTeacherHandler)
	r.GET("/set/teacher/stavka", AuthRequired(), directory.GetSetTeacherStavkaHandler)
	r.GET("/set/teacher/paiment", AuthRequired(), directory.GetSetTeacherPaimentHandler)
	r.DELETE("/delete/teacher", AuthRequired(), directory.GetDirectoryDeleteTeacherHandler)

	// currency
	r.GET("/list/currency", AuthRequired(), directory.GetDirectoryListCurrencyHandler)
	r.POST("/input/currency", AuthRequired(), directory.PostDirectoryInputCurrencyHandler)
	r.GET("/view/currency", AuthRequired(), directory.GetDirectoryViewCurrencyHandler)
	r.DELETE("/delete/currency", AuthRequired(), directory.GetDirectoryDeleteCurrencyHandler)

	port := config.C.UrlPort

	r.Run(fmt.Sprintf(":%d", port))
}

func startPageHandler(c *gin.Context) {
	user := access.GetUser(c)
	if c.IsAborted() {
		return
	}
	entrance.RedirectStartPageByUserRole(c, user.Role)
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func helpHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "help.html", nil)
}

func aboutHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", nil)
}

func AuthRequired() gin.HandlerFunc {

	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			c.Redirect(http.StatusFound, "/index")
			c.Abort()
			return
		}
		c.Next()
	}
}
