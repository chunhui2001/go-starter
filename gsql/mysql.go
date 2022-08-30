package gsql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chunhui2001/go-starter/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var (
	dbClient *sql.DB
	logger   *logrus.Entry
)

func Init(log *logrus.Entry) {

	logger = log

	connectionString := "root:Cc@tcp(localhost:3306)/mydb?timeout=90s&interpolateParams=true"
	// connectionString := "root:Cc@tcp(localhost:3306)/mydb?timeout=90s"

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		panic(err)
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()

	if err != nil {
		logger.Error(fmt.Sprintf("Mysql-Client-Connect-Error: connectionString=%s, errorMessage=%s", connectionString, utils.ErrorToString(err)))
		return
	}

	dbClient = db

	logger.Info(fmt.Sprintf("Mysql-Client-Connected-Successful"))

	CreateTable(`create table if not exists book(isbn varchar(14), title varchar(200), price int, primary key(isbn))`)
	Insert(`insert into book(isbn, title, price) values(?, ?, ?)`, "978-4798161490", "MySQL徹底入門 第4版", 4180)
	Insert(`insert into book(isbn, title, price) values(?, ?, ?)`, "978-4798161491", "MySQL徹底入門 第4版", 4180)
	Insert(`insert into book(isbn, title, price) values(?, ?, ?)`, "978-4798161492", "MySQL徹底入門 第4版", 4180)

}

func CreateTable(ddl string) {

	_, err := dbClient.Exec(ddl)

	if err != nil {
		logger.Error(fmt.Sprintf("Mysql-Create-Table-Error: ddl=%s, errorMessage=%s", ddl, utils.ErrorToString(err)))
	}

	logger.Info(fmt.Sprintf("Mysql-Create-Table-Successful"))

}

func Insert(sql string, args ...any) sql.Result {

	result, err := dbClient.Exec(sql, args...)

	if err != nil {
		logger.Error(fmt.Sprintf("Mysql-Insert-Error: sql=%s, errorMessage=%s", sql, string(err.Error())))
		return nil
	}

	return result

}
