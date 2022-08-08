package logger

import (
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/utils"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

func init() {

	env := config.GetEnv("GIN_ENV", "development")
	app := config.GetEnv("APP_NAME", "go-starter")

	log_folder := "/tmp/" + app

	lumberjackLogger := &lumberjack.Logger{
		Filename:   log_folder + "/mylog.txt",
		MaxSize:    1, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
		Compress:   true,
	}

	// assign it to the standard logger
	mw := io.MultiWriter(os.Stdout, lumberjackLogger)

	Log = logrus.New()

	Log.SetOutput(mw)
	Log.SetLevel(logrus.DebugLevel)
	Log.SetReportCaller(true)

	Log.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	Log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output

	Log.SetFormatter(&MyFormatter{
		TimestampFormat: utils.TimeStampFormat,
		LogFormat:       "%time% [%lvl%] - %file% >> %msg%\n",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			//return frame.Function, fileName
			return "", fileName
		},
	})

	// config gin
	if env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DisableConsoleColor() //  Disable Console Color
	gin.DefaultWriter = io.MultiWriter(lumberjackLogger, os.Stdout)

	Log.WithFields(logrus.Fields{
		"App": app,
		"Env": env,
	}).Info("Initialization log completed: app=", app, ", env=", env)

	// don't forget to close it
	//defer f.Close()

}
