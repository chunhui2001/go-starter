package gsql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/core/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type MySql struct {
	Enable       bool   `mapstructure:"MYSQL_ENABLE"`
	Dns          string `mapstructure:"MYSQL_DNS"`
	InitScript   string `mapstructure:"MYSQL_INIT_SCRIPT"`
	UpdateScript string `mapstructure:"MYSQL_UPDATE_SCRIPT"`
}

var (
	dbClient  *sql.DB
	logger    *logrus.Entry
	mySqlConf *MySql
)

func Init(conf *MySql, log *logrus.Entry) {

	logger = log
	mySqlConf = conf

	hostMatch := utils.Matches(mySqlConf.Dns, `\(([0-9\.a-zA-Z_]+:[0-9]+)?\)/([A-Za-z0-9_]+)`)
	hostName := hostMatch[0][1] + "/" + hostMatch[0][2]
	db, err := sql.Open("mysql", mySqlConf.Dns)

	if err != nil {
		panic(err)
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	dbClient = db

	if err := dbClient.Ping(); err != nil {
		logger.Error(fmt.Sprintf("Mysql-Client-Connect-Error: MySqlServer=%s, errorMessage=%s", hostName, string(err.Error())))
		return
	}

	if version, err := Version(); err == nil {

		logger.Info(fmt.Sprintf("Mysql-Client-Connected-Successful: MySqlServer=%s, Version=%s", hostName, version))

		// execute the Embedding scripts
		exceScripts()

		return

	}

	logger.Error(fmt.Sprintf("Mysql-Client-Connect-Error: MySqlServer=%s, errorMessage=%s", hostName, string(err.Error())))

}

func Version() (string, error) {
	var version string
	err2 := dbClient.QueryRow("SELECT VERSION()").Scan(&version)
	return version, err2
}

func exceScripts() (bool, error) {

	if mySqlConf.InitScript != "" {

		initScriptFolder := filepath.Join(utils.RootDir(), mySqlConf.InitScript)

		if ok, _ := utils.FileExists(initScriptFolder); !ok {
			logger.Warnf("MySql-InitScript-Folder-Not-Exists: InitScriptFolder=%s", initScriptFolder)
		} else {

			fileItemArray, keys := getSortedFiles(initScriptFolder)

			for _, number := range keys {

				item := fileItemArray[number].(map[string]interface{})
				status := ""

				if fileByte, err := utils.ReadFile(item["path"].(string)); err == nil {
					if len(fileByte) > 0 {
						if _, fail := ExecuteDdlScripts(string(fileByte)); fail != nil {
							status = "Executed-Error(" + fail.Error() + ")"
						} else {
							status = "Executed-Completed"
						}
					} else {
						status = "Empty-File"
					}
				} else {
					status = "Read-Error(" + err.Error() + ")"
				}

				logger.Infof("Mysql-IniterScript-File-Loading: status=%s, path=%s/%s", status, mySqlConf.InitScript, item["name"])

			}

		}

	}

	if mySqlConf.UpdateScript != "" {

		updateScriptFolder := filepath.Join(utils.RootDir(), mySqlConf.UpdateScript)

		if ok, _ := utils.FileExists(updateScriptFolder); !ok {
			logger.Warnf("MySql-UpdateScript-Folder-Not-Exists: InitScriptFolder=%s", updateScriptFolder)
		} else {

			fileItemArray, keys := getSortedFiles(updateScriptFolder)

			for _, number := range keys {

				item := fileItemArray[number].(map[string]interface{})
				status := ""

				if fileByte, err := utils.ReadFile(item["path"].(string)); err == nil {
					if len(fileByte) > 0 {
						if _, fail := ExecuteDdlScripts(string(fileByte)); fail != nil {
							status = "Executed-Error(" + fail.Error() + ")"
						} else {
							status = "Executed-Completed"
						}
					} else {
						status = "Empty-File"
					}
				} else {
					status = "Read-Error(" + err.Error() + ")"
				}

				logger.Infof("Mysql-UpdateScript-File-Loading: status=%s, path=%s/%s", status, mySqlConf.UpdateScript, item["name"])

			}

		}

	}

	return true, nil

}

func ExecuteDdlScripts(ddl string) (bool, error) {

	_, err := dbClient.Exec(ddl)

	if err != nil {
		return false, err
	}

	return true, nil

}

// ss := &SimpleSelect{
// 			Table:  "t_books",
// 			Fields: []string{"f_id", "f_title", "f_created_at"},
// 			Params: utils.MapOf("f_id", 1, "f_title", "sd"),
// 		}

// https://www.golang-book.com/books/intro/8
func SimpleQuery(selects *SimpleSelect) ([]map[string]interface{}, error) {

	xselect, valus := selects.ToString()
	rows, err := dbClient.Query(xselect, valus...)

	if err != nil {
		logger.Infof("Mysql-SimpleQuery-Error: sql=%s, valus=%s, error=%s", xselect, utils.ToJsonString(valus), err.Error())
		return nil, err
	}

	defer rows.Close()

	cols, err := rows.Columns()

	if cols == nil {
		return nil, err
	}

	colTypes, err2 := rows.ColumnTypes()

	if err2 != nil {
		return nil, err2
	}

	var result []map[string]interface{}

	for rows.Next() {

		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		values := make([]interface{}, len(cols))

		for i := range values {
			currType := colTypes[i].DatabaseTypeName()
			if currType == "INT" {
				values[i] = new(int32)
			} else if currType == "VARCHAR" {
				values[i] = new(string)
			} else if currType == "TIMESTAMP" {
				values[i] = new(string)
			} else {
				logger.Errorf("Mysql-Current-Data-Type-Not-Cached: DatabaseTypeName=%s, ColumeName=%s", currType, colTypes[i].Name())
				values[i] = new(interface{})
			}
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		currentRow := make(map[string]interface{})

		for i, colName := range cols {
			currentRow[colName] = values[i]
		}

		result = append(result, currentRow)
	}

	return result, nil

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

// ss := &SimpleSelect{
// 			Table:  "t_books",
// 			Fields: []string{"f_id", "f_title", "f_created_at"},
// 			Params: utils.MapOf("f_id", 1, "f_title", "sd"),
// 		}

// 		xs, vals := ss.ToString()

type SimpleSelect struct {
	Table  string
	Fields []string
	Params map[string]interface{}
}

func (s *SimpleSelect) ToString() (string, []any) {

	fieldsString := "*"

	if len(s.Fields) > 0 {
		fieldsString = "`" + strings.Join(s.Fields, "`, `") + "`"
	}

	flds := []string{}
	vals := []any{}

	for k, v := range s.Params {
		flds = append(flds, "`"+k+"` = ?")
		vals = append(vals, v)
	}

	pl := strings.Join(flds, " and ")

	if len(flds) == 0 {
		pl = "1=1"
	}

	xselect := fmt.Sprintf("select %s from `%s` where %s;", fieldsString, s.Table, pl)

	return xselect, vals

}

func getSortedFiles(file_path string) (map[int]interface{}, []int) {

	files, _ := ioutil.ReadDir(file_path)
	var fileItemArray = make(map[int]interface{})

	for _, file := range files {

		fileNumberMatch := utils.Matches(file.Name(), `^([0-9]+)_((.)*\.sql)$`)

		if fileNumberMatch != nil {
			filenumber := fileNumberMatch[0][1]
			if number, err := strconv.Atoi(filenumber); err == nil {
				fileItemArray[number] = utils.MapOf(
					"number", number,
					"name", file.Name(),
					"path", filepath.Join(file_path, file.Name()),
					"scripts", "",
				)
			}
		}

	}

	fileItemArray, keys := utils.SortedKeysInt(fileItemArray)

	return fileItemArray, keys

}
