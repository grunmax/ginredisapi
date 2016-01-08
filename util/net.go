package util

import (
	"github.com/gin-gonic/gin"
)

func TestKeys() []string {
	return []string{"alfa", "zulu", "charlie"}
}

func BodyOk(message string) map[string]string {
	result := map[string]string{}
	result["status"] = "OK"
	result["message"] = message
	return result
}

func BodyErr(message string) map[string]string {
	result := map[string]string{}
	result["status"] = "ERROR"
	result["message"] = message
	return result
}

func GetCookieValue(name string, c *gin.Context) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	} else {
		return cookie.Value, nil
	}

}
