package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	_ "strings"
	"time"

	"io"
	"path"
	"path/filepath"
	"runtime"

	"github.com/chunhui2001/go-starter/gredis"
	"github.com/chunhui2001/go-starter/utils"
	_ "github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"

	lkh "github.com/gfremex/logrus-kafka-hook"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().Format(utils.TimeStampFormat) + " [STDOUT] - " + string(bytes))
}

type App struct {
	Env     string `mapstructure:"GIN_ENV"`
	AppName string `mapstructure:"APP_NAME"`
	AppPort string `mapstructure:"APP_PORT"`
	LogFile string `mapstructure:"APP_LOG_FILE"`
}

type Wss struct {
	Enable bool   `mapstructure:"WSS_ENABLE"`
	Prefix string `mapstructure:"WSS_PREFIX"`
	Host   string `mapstructure:"WSS_HOST"`
}

type Cookie struct {
	Enable bool   `mapstructure:"COOKIE_ENABLE"`
	Name   string `mapstructure:"COOKIE_NAME"`
	Secret string `mapstructure:"COOKIE_SECRET"`
	MaxAge int    `mapstructure:"COOKIE_MaxAge"`
}

func (w *Wss) Wss() string {
	return w.Host + w.Prefix
}

var AppSetting = &App{
	// Env:     "development",
	Env:     "production",
	AppName: "go-starter",
	AppPort: "8080",
	LogFile: "",
}

var WssSetting = &Wss{
	Enable: false,
	Prefix: "",
}

var CookieSetting = &Cookie{
	Enable: false,
	Name:   "_mycookie_name",
	Secret: "_mycookie_secret",
	MaxAge: 1 * 60,
}

var RedisConf = &gredis.GRedis{
	Mode:           gredis.Disabled,
	MasterName:     "",
	Host:           "127.0.0.1:6379",
	Addrs:          "",
	Db:             -1,
	Passwd:         "",
	RouteByLatency: false,
	RouteRandomly:  false,
}

var Log *logrus.Entry
var filename string = ".env"

// LoadEnvVars will load a ".env[.development|.test]" file if it exists and set ENV vars.
// Useful in development and test modes. Not used in production.
func init() {

	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	var env string = os.Getenv("GIN_ENV")
	var envfile = ".env." + env

	if exists, _ := utils.FileExists(filepath.Join(utils.RootDir(), envfile)); exists == true {
		filename = envfile
	} else {
		log.Println("Configuration loading " + filepath.Join(utils.RootDir(), envfile) + " file error, use .env file.")
	}

	v1 := readConfig(filename, map[string]interface{}{})

	loadAppSettings(v1, filename)

	// init log configuration
	InitLog()

	loadWssSettings(v1, filename)
	loadRedisSettings(v1, filename)
	loadCookieSettings(v1, filename)

}

type CallInfo struct {
	packageName string
	fileName    string
	funcName    string
	line        int
}

func InitLog() {

	env := AppSetting.Env
	app := AppSetting.AppName

	log_file := filepath.Join(utils.TempDir(), app, "mylog.txt")

	lumberjackLogger := &lumberjack.Logger{
		Filename:   log_file,
		MaxSize:    1, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
		Compress:   true,
	}

	// assign it to the standard logger
	mw := io.MultiWriter(os.Stdout, lumberjackLogger)

	myLog := logrus.New()

	myLog.SetOutput(mw)
	myLog.SetLevel(logrus.DebugLevel)
	myLog.SetReportCaller(true)

	myLog.SetFormatter(&MyTxtFormatter{
		TimestampFormat: utils.TimeStampFormat,
		LogFormat:       "%time% [%lvl%] - %file% >> %msg%\n",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return "", fmt.Sprintf("%s() %s:%d", frame.Function, path.Base(frame.File), frame.Line)
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

	kafkaServerAddr := "127.0.0.1:9092"
	kafkaServer := []string{kafkaServerAddr}

	hook, err := lkh.NewKafkaHook(
		"kh",
		[]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel},
		// &logrus.JSONFormatter{},
		&MyJSONFormatter{},
		kafkaServer,
	)

	if err != nil {
		myLog.Error(fmt.Sprintf("Kafka Append Initialization failed: kafkaServer=%s, errorMessage=%s", kafkaServerAddr, utils.ErrorToString(err)))
	}

	myLog.Hooks.Add(hook)

	// logrus.WithField("app_name", []string{"go-starter"})
	Log = myLog.WithField("topics", []string{"topic_1"})

	Log.Info("Initialization log completed: appRoot=", utils.RootDir(), ", logFile=", log_file)

}

func loadAppSettings(v1 *viper.Viper, filename string) {
	err := v1.Unmarshal(&AppSetting)
	if err != nil {
		log.Println("viper parse AppSettings error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		log.Println("AppSetting: Env=" + AppSetting.Env + ", AppName=" + AppSetting.AppName + ", AppPort=" + AppSetting.AppPort)
	}
}

func loadRedisSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&RedisConf)

	if err != nil {
		Log.Info("viper parse RedisConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		if RedisConf.Mode == 0 {
			Log.Info("RedisSettings: Mode=" + utils.ToString(RedisConf.Mode) + ", Host=" + RedisConf.Host)
		} else {
			gredis.Init(RedisConf, Log)
		}
	}

}

func loadWssSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&WssSetting)

	if err != nil {
		Log.Info("viper parse WssSetting error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		Log.Info("WssSetting: Host=" + WssSetting.Host + ", Prefix=" + WssSetting.Prefix)
	}

}

func loadCookieSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&CookieSetting)

	if err != nil {
		Log.Info("viper parse CookieSetting error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		Log.Info("CookieSetting: Enable=" + strconv.FormatBool(CookieSetting.Enable) + ", Name=" + CookieSetting.Name + ", Secret=" + CookieSetting.Secret + ", MaxAge=" + fmt.Sprint(CookieSetting.MaxAge))
	}

}

func readConfig(filename string, defaults map[string]interface{}) *viper.Viper {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.AddConfigPath(utils.RootDir())
	v.SetConfigName(filename)
	v.SetConfigType("env")
	// v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		log.Println("viper loaded error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return nil
	}
	log.Println("viper Configuration loaded " + filepath.Join(utils.RootDir(), filename) + " successful.")
	return v
}

// GetEnv returns an environment variable or a default value if not present
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
