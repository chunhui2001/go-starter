package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"content": "This is an Home page...",
	})
}

func AboutRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "about", gin.H{
		"content": "This is an about page...",
	})
}
