package logger

import (
	"errors"
	// "fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"

	"go-starter/config"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var Log *logrus.Logger
var timeStampFormat = "2006-01-02T15:04:05.000Z07:00"

func init() {

	env := config.GetEnv("GIN_ENV", "development")
	app := config.GetEnv("APP_NAME", "go-starter")

	log_folder := "/tmp/" + app

	if _, err := os.Stat(log_folder); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(log_folder, os.ModePerm)
	}

	// open a file
	f, _ := os.OpenFile(log_folder+"/log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	// assign it to the standard logger
	mw := io.MultiWriter(os.Stdout, f)

	Log = logrus.New()

	Log.SetOutput(mw)
	Log.SetLevel(logrus.DebugLevel)

	Log.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	Log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output

	Log.SetFormatter(&easy.Formatter{
		TimestampFormat: timeStampFormat,
		LogFormat:       "[%lvl%] %time% - %msg%\n",
	})

	// config gin
	if env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DisableConsoleColor() //  Disable Console Color
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	Log.WithFields(logrus.Fields{
		"App": app,
		"Env": env,
	}).Info("Initialization log completed: app=", app, ", env=", env)

	// don't forget to close it
	//defer f.Close()

}
