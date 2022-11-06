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
	Opts         string `mapstructure:"MYSQL_CONN_OPTS"`
	Server       string `mapstructure:"MYSQL_SERVER"`
	Database     string `mapstructure:"MYSQL_DATABASE"`
	User         string `mapstructure:"MYSQL_USER_NAME"`
	Passwd       string `mapstructure:"MYSQL_PASSWD"`
	InitScript   string `mapstructure:"MYSQL_INIT_SCRIPT"`
	UpdateScript string `mapstructure:"MYSQL_UPDATE_SCRIPT"`
}

func (c *MySql) connString(passwd string) string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?%s`, c.User, passwd, c.Server, c.Database, c.Opts)
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
		return
	}

	if version, err := Version(); err == nil {
		logger.Info(fmt.Sprintf("Mysql-Client-Connected-Successful: ServerVersion=%s, ConnString=%s", version, conf.connString("****")))
		// execute the Embedding scripts
		exceScripts()
		return

	}

	logger.Error(fmt.Sprintf("Mysql-Client-Connect-Error: ConnString=%s, errorMessage=%s", conf.connString("****"), string(err.Error())))

}

func Version() (string, error) {
	var version string
	err2 := DbClient.QueryRow("SELECT VERSION()").Scan(&version)
	return version, err2
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
func Exec(sqlStr string, args ...any) sql.Result {

	// prepare the statement
	stmt, _ := DbClient.Prepare(sqlStr)

	// format all args at once
	result, err := stmt.Exec(args...)

	if err != nil {
		logger.Error(fmt.Sprintf("Mysql-Insert-Error: sqlStr=%s, errorMessage=%s", sqlStr, string(err.Error())))
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
