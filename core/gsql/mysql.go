package gsql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chunhui2001/go-starter/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type MySql struct {
	Enable bool   `mapstructure:"MYSQL_ENABLE"`
	Dns    string `mapstructure:"MYSQL_DNS"`
}

var (
	dbClient *sql.DB
	logger   *logrus.Entry
)

func Init(mySqlConf *MySql, log *logrus.Entry) {

	logger = log
	dns := mySqlConf.Dns
	hostMatch := utils.Matches(dns, `\(([0-9\.a-zA-Z_]+:[0-9]+)?\)/([A-Za-z0-9_]+)`)
	hostName := hostMatch[0][1] + "/" + hostMatch[0][2]
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	dbClient = db

	err = dbClient.Ping()

	if err != nil {
		logger.Error(fmt.Sprintf("Mysql-Client-Connect-Error: MySqlServer=%s, errorMessage=%s", hostName, string(err.Error())))
		return
	}

	logger.Info(fmt.Sprintf("Mysql-Client-Connected-Successful: MySqlServer=%s", hostName))

}

// CreateTable(`create table if not exists book(isbn varchar(14), title varchar(200), price int, primary key(isbn))`)
func CreateTable(ddl string) {

	_, err := dbClient.Exec(ddl)

	if err != nil {
		logger.Error(fmt.Sprintf("Mysql-Create-Table-Error: ddl=%s, errorMessage=%s", ddl, string(err.Error())))
	}

	logger.Info(fmt.Sprintf("Mysql-Create-Table-Successful"))

}

// sqlStr := "INSERT INTO test(n1, n2, n3) VALUES "
// vals := []interface{}{}

// for _, row := range data {
//    sqlStr += "(?, ?, ?),"
//    vals = append(vals, row["v1"], row["v2"], row["v3"])
// }

// //trim the last ,
// sqlStr = strings.TrimSuffix(sqlStr, ",")

// //Replacing ? with $n for postgres
// sqlStr = ReplaceSQL(sqlStr, "?")

// //prepare the statement
// stmt, _ := db.Prepare(sqlStr)

// //format all vals at once
// res, _ := stmt.Exec(vals...)
func Exec(sqlStr string, args ...any) sql.Result {

	// prepare the statement
	stmt, _ := dbClient.Prepare(sqlStr)

	// format all args at once
	result, err := stmt.Exec(args...)

	if err != nil {
		logger.Error(fmt.Sprintf("Mysql-Insert-Error: sqlStr=%s, errorMessage=%s", sqlStr, string(err.Error())))
		return nil
	}

	return result

}

// Insert(`insert into book(isbn, title, price) values(?, ?, ?)`, "978-4798161495", "MySQL徹底入門 第4版", 4180)
func Insert(sql string, args ...any) sql.Result {

	result, err := dbClient.Exec(sql, args...)

	if err != nil {
		logger.Error(fmt.Sprintf("Mysql-Insert-Error: sql=%s, errorMessage=%s", sql, string(err.Error())))
		return nil
	}

	return result

}
