package gsql

import (
	"context"
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
	"github.com/thoas/go-funk"
)

type MySql struct {
	Enable       bool   `mapstructure:"MYSQL_ENABLE"`
	Opts         string `mapstructure:"MYSQL_CONN_OPTS" json:"opts"`
	Server       string `mapstructure:"MYSQL_SERVER" json:"server"`
	Database     string `mapstructure:"MYSQL_DATABASE" json:"database"`
	User         string `mapstructure:"MYSQL_USER_NAME" json:"user_name"`
	Passwd       string `mapstructure:"MYSQL_PASSWD" json:"passwd"`
	InitScript   string `mapstructure:"MYSQL_INIT_SCRIPT"`
	UpdateScript string `mapstructure:"MYSQL_UPDATE_SCRIPT"`
}

func (c *MySql) connString(passwd string) string {
	// return fmt.Sprintf(`%s:%s@tcp(%s)/%s?%s`, c.User, passwd, c.Server, c.Database, c.Opts)
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s`, c.User, passwd, c.Server, c.Database)
}

var (
	DbClient  *sql.DB
	logger    *logrus.Entry
	mySqlConf *MySql
)

func Init(conf *MySql, log *logrus.Entry) {

	logger = log
	mySqlConf = conf

	db, err := sql.Open("mysql", conf.connString(conf.Passwd))

	if err != nil {
		panic(err)
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	DbClient = db

	if err := DbClient.Ping(); err != nil {
		logger.Error(fmt.Sprintf("Mysql-Client-Connect-Error: ConnString=%s, errorMessage=%s", conf.connString("****"), string(err.Error())))
		panic(err)
	}

	if version, err := Version(); err == nil {
		logger.Info(fmt.Sprintf("Mysql-Client-Connected-Successful: ServerVersion=%s, ConnString=%s", version, conf.connString("****")))
		// execute the Embedding scripts
		exceScripts()
		return

	}

	logger.Error(fmt.Sprintf("Mysql-Client-Connect-Error: ConnString=%s, errorMessage=%s", conf.connString("****"), string(err.Error())))

}

func Client() *sql.DB {
	return DbClient
}

func Version() (string, error) {
	var version string
	err2 := DbClient.QueryRow("SELECT VERSION()").Scan(&version)
	return version, err2
}

func ShowTables() (string, error) {
	var talbes string
	err2 := DbClient.QueryRow("show tables;").Scan(&talbes)
	return talbes, err2
}

func exceScripts() {

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

}

func ExecuteDdlScripts(ddl string) (bool, error) {

	_, err := DbClient.Exec(ddl)

	if err != nil {
		return false, err
	}

	// logger.Infof("Mysql-ExecuteDdlScripts: ddl=%s, IsError=%t", ddl, err != nil)

	return true, nil

}

func StatmentQuery(sql string) {

	// DbClient.Exec("")

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
func Exec(sqlStr string, timeout int32, args ...any) (sql.Result, error) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Nanosecond*1000000000*time.Duration(timeout))
	defer cancelFunc()

	// prepare the statement
	stmt, err := DbClient.Prepare(sqlStr)

	if err != nil {
		theSql := sqlStr
		if len(sqlStr) > 135 {
			theSql = "`" + theSql[0:135] + "` ...."
		}
		logger.Error(fmt.Sprintf("Mysql-Insert-Error-1: sqlStr=%s, server=%s, errorMessage=%s", theSql, mySqlConf.connString("****"), string(err.Error())))
		return nil, err
	}

	// format all args at once
	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		logger.Error(fmt.Sprintf("Mysql-Insert-Error-2: sqlStr=%s, server=%s, errorMessage=%s", sqlStr, mySqlConf.connString("****"), string(err.Error())))
		return nil, err
	}

	return result, nil

}

// 批量插入
// INSERT INTO tableName(id, name) VALUES (?, ?),(?, ?),(?, ?)
func InsertBulk(timeout int32, tableName string, columeMaps [][]string, insertData []map[string]interface{}) (sql.Result, error) {

	_columes := funk.Map(columeMaps, func(m []string) string { return m[0] })
	_keys := funk.Map(columeMaps, func(m []string) string { return m[1] })

	insert := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", tableName, strings.Join(_columes.([]string)[:], ", "))

	vals := make([]any, 0, len(insertData))
	placeholders := strings.Repeat("?, ", len(_columes.([]string)))

	for _, item := range insertData {
		insert += fmt.Sprintf(`(%s),`, placeholders[:len(placeholders)-2])
		for _, _k := range _keys.([]string) {
			if item[_k] == nil {
				vals = append(vals, nil)
			} else {
				vals = append(vals, item[_k].(interface{}))
			}
		}
	}

	return Exec(insert[:len(insert)-1]+";", timeout, vals...)

}

// ss := &SimpleSelect{
// 			Table:  "t_books",
// 			Fields: []string{"f_id", "f_title", "f_created_at"},
// 			Params: utils.MapOf("f_id", 1, "f_title", "sd"),
// 		}

// 		xs, vals := ss.ToString()

type SimpleSelect struct {
	Table     string
	Fields    []string
	Params    map[string]interface{}
	BeginTrx  bool
	ForUpdate bool
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

	beginTrx := ""

	// if s.BeginTrx {
	// 	beginTrx = "begin; "
	// }

	forUpdate := ""

	if s.ForUpdate {
		forUpdate = " for update"
	}

	xselect := fmt.Sprintf("%sselect %s from `%s` where %s%s;", beginTrx, fieldsString, s.Table, pl, forUpdate)

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
