package actions

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/gsql"
	_ "github.com/chunhui2001/go-starter/core/utils"

	"github.com/gin-gonic/gin"
)

var (
	DbClient  *sql.DB
	mySqlConf *gsql.MySql
)

func connString(user string, passwd string, server string, db string, opts string) string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?%s`, user, passwd, server, db, opts)
}

func Version() (string, error) {
	var version string
	err2 := DbClient.QueryRow("SELECT VERSION()").Scan(&version)
	return version, err2
}

func ListTables(c *gin.Context) {

	var mysql = &gsql.MySql{}

	if err := c.ShouldBindJSON(&mysql); err != nil {
		c.JSON(200, (&R{Error: err}).Fail(413))
		return
	}

	var connString = connString(mysql.User, mysql.Passwd, mysql.Server, mysql.Database, mysql.Opts)

	db, err := sql.Open("mysql", connString)

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(500))
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	DbClient = db

	if err := DbClient.Ping(); err != nil {
		c.JSON(200, (&R{Error: err}).Fail(500))
		return
	}

	if version, err := Version(); err == nil {
		logger.Info(fmt.Sprintf("Mysql-Client-Connected-Successful: ServerVersion=%s, ConnString=%s", version, connString))
	} else {
		logger.Error(fmt.Sprintf("Mysql-Client-Connect-Error: ConnString=%s, errorMessage=%s", connString, string(err.Error())))
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Nanosecond*1000000*100)

	defer cancelFunc()

	var rows *sql.Rows
	rows, err = DbClient.QueryContext(ctx, "show tables")

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(500))
		return
	}

	tables := make([]string, 0)

	for rows.Next() {
		var talbeName string
		if err2 := rows.Scan(&talbeName); err2 != nil {
			c.JSON(200, (&R{Error: err2}).Fail(500))
			return
		}
		tables = append(tables, talbeName)
	}

	c.JSON(200, (&R{Data: tables}).Success())

}
