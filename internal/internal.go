package internal

import (
	"fmt"
	"schoolonline/webmessage"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Helper function to check required fields
func CheckRequiredFields(fields ...string) string {
	fieldNames := []string{"school_id", "student_id", "parent_id", "currency_id"}
	for i, field := range fields {
		if field == "" {
			return fieldNames[i]
		}
	}
	return ""
}

func GetQueryString(c *gin.Context, key string) (res string) {
	res = c.Query(key)
	if res == "" {
		webmessage.Err(c, fmt.Errorf("%s not provided", key), "data не указан", "/menu")
		c.Abort()
		return
	}
	return
}

func GetFormaString(c *gin.Context, key string) (res string) {
	res = c.PostForm(key)
	if res == "" {
		webmessage.Err(c, fmt.Errorf("%s not provided", key), "data не указан", "/menu")
		c.Abort()
		return
	}
	return
}

func GetFormaInt(c *gin.Context, key string) (res int) {
	rs := c.PostForm(key)
	if rs == "" {
		webmessage.Err(c, fmt.Errorf("%s not provided", key), "data не указан", "/menu")
		c.Abort()
		return
	}

	res, err := strconv.Atoi(rs)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ID не указан", "/menu")
		c.Abort()
		return
	}
	return
}

func GetFormaFloat(c *gin.Context, key string) (res float64) {
	rs := c.PostForm(key)
	if rs == "" {
		webmessage.Err(c, fmt.Errorf("%s not provided", key), "data не указан", "/menu")
		c.Abort()
		return
	}

	res, err := strconv.ParseFloat(rs, 64)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "ID не указан", "/menu")
		c.Abort()
		return
	}
	return
}

func GetFormaDate(c *gin.Context, key string) (res time.Time) {
	rs := c.PostForm(key)
	if rs == "" {
		webmessage.Err(c, fmt.Errorf("%s not provided", key), "date не указан", "/menu")
		c.Abort()
		return
	}

	res, err := time.Parse("2006-01-02", rs)
	if err != nil {
		fmt.Println(err)
		webmessage.Err(c, err, "date не указан", "/menu")
		c.Abort()
		return
	}
	return
}
