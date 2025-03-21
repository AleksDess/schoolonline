package user

import (
	"fmt"
	"net/http"
	"net/url"
	"schoolonline/config"
	"schoolonline/dict"
	"schoolonline/internal"
	"schoolonline/launch"
	"schoolonline/qrcode"
	"schoolonline/webmessage"

	"github.com/gin-gonic/gin"
)

type LinkTgBot struct {
	Link string `db:"link" json:"link" form:"link"`
}

func GetUserProfileHandler(c *gin.Context) {

	login := internal.GetQueryString(c, "login")

	us, err := dict.GetUserByLogin(login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "", "/")
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "user_profile.html", gin.H{
		"data": us,
	})

}

func GetUserSettingHandler(c *gin.Context) {

	login := internal.GetQueryString(c, "login")

	us, err := dict.GetUserByLogin(login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "", "/")
		c.Abort()
		return
	}

	link := LinkTgBot{Link: config.C.LinkTest}

	if launch.Launch == "server" {
		link.Link = config.C.LinkServer
	}

	linkQrCode := link.Link + us.Login

	qrBase64 := qrcode.CreateQrCode(linkQrCode)

	c.HTML(http.StatusOK, "user_setting.html", gin.H{
		"data":   us,
		"link":   link,
		"QRCode": qrBase64,
	})

}

func PostUserProfileHandler(c *gin.Context) {

	login := internal.GetFormaString(c, "fullName")

	us, err := dict.GetUserByLogin(login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "", "/")
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "user_profile.html", gin.H{
		"data": us,
	})

}

func PostUserSettingHandler(c *gin.Context) {

	login := internal.GetFormaString(c, "fullName")

	us, err := dict.GetUserByLogin(login)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "", "/")
		c.Abort()
		return
	}

	link := LinkTgBot{Link: config.C.LinkTest}

	if launch.Launch == "server" {
		link.Link = config.C.LinkServer
	}

	linkQrCode := link.Link + us.Login

	qrBase64 := qrcode.CreateQrCode(linkQrCode)

	c.HTML(http.StatusOK, "user_setting.html", gin.H{
		"data":   us,
		"link":   link,
		"QRCode": qrBase64,
	})

}

func PostUserEditEmailHandler(c *gin.Context) {
	login := internal.GetFormaString(c, "edit-email")
	if c.IsAborted() {
		return
	}

	c.Redirect(http.StatusFound, "/user/profile?login="+url.QueryEscape(login))
}

func PostUserEditTgHandler(c *gin.Context) {
	login := internal.GetFormaString(c, "edit-tg")
	if c.IsAborted() {
		return
	}

	c.Redirect(http.StatusFound, "/user/profile?login="+url.QueryEscape(login))
}

func GetCheckLoginHandler(c *gin.Context) {

	login := c.Query("login")
	if login == "" {
		webmessage.Err(c, fmt.Errorf("id not provided"), "ID факультета не указан", "/directory")
		c.Abort()
		return
	}

	fmt.Println(login)

	exists, _, _ := dict.UserExists(login)
	c.JSON(http.StatusOK, gin.H{"exists": exists})
}
