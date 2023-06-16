package actions

import (
	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/gsql"
	"github.com/chunhui2001/go-starter/core/utils"

	"github.com/gin-gonic/gin"
)

func QueryBooksRouter(c *gin.Context) {

	ss := &gsql.SimpleSelect{
		Table:  "t_books",
		Fields: []string{"f_id", "f_title", "f_created_at"},
		Params: utils.MapOf("f_title", "MySQL徹底入門 第4版"),
	}

	rows, err := gsql.SimpleQuery(ss)

	c.JSON(200, (&R{Data: rows, Error: err}).IfErr(400))

}
