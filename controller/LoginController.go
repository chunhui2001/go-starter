package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type PostLoginForm struct {
	Username string `form:"username"`
	Passwd   string `form:"passwd"`
}

var users = map[string]string{
	"user1": "1111",
	"user2": "2222",
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Header("Cache-Control", "no-cache")
	c.Redirect(302, "/")
}

func LoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"loginUrl": WEB_PAGE_CONF.LoginUrl,
	})
}

func PostLoginHandler(c *gin.Context) {

	var formData PostLoginForm

	if err := c.ShouldBindWith(&formData, binding.FormMultipart); err != nil {
		c.HTML(http.StatusOK, "login", gin.H{
			"loginUrl":               WEB_PAGE_CONF.LoginUrl,
			"fromMessageForUserName": "This is a required field.",
		})
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[formData.Username]

	if !ok || expectedPassword != formData.Passwd {
		c.HTML(http.StatusOK, "login", gin.H{
			"loginUrl":               WEB_PAGE_CONF.LoginUrl,
			"fromMessageForUserName": "Invalid Username Or Password." + formData.Passwd + ", " + formData.Username,
		})
		return
	}

	session := sessions.Default(c)
	session.Set("username", formData.Username)
	defer session.Save()

	c.Redirect(http.StatusFound, "/")

}

func SignUpHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "signup", gin.H{
		"content": "This is signup page...",
	})
}
