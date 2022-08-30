package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"io"
	"path"
	"path/filepath"
	"runtime"

	"github.com/chunhui2001/go-starter/gmongo"
	"github.com/chunhui2001/go-starter/gredis"
	"github.com/chunhui2001/go-starter/utils"
	_ "github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"

	_ "github.com/chunhui2001/go-starter/gmongo"
	lkh "github.com/gfremex/logrus-kafka-hook"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var timeStampFormat = "2006-01-02T15:04:05Z07:00"

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().Format(timeStampFormat) + " [STDOUT] - " + string(bytes))
}

type App struct {
	Env      string `mapstructure:"GIN_ENV"`
	AppName  string `mapstructure:"APP_NAME"`
	AppPort  string `mapstructure:"APP_PORT"`
	TimeZone string `mapstructure:"APP_TIMEZONE"`
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

type LogConf struct {
	Output      string `mapstructure:"LOG_OUTPUT"`
	FilePath    string `mapstructure:"LOG_FILE_PATH"`
	KafkaServer string `mapstructure:"LOG_KAFKA_SERVER"`
	KafkaTopic  string `mapstructure:"LOG_KAFKA_TOPIC"`
}

func (w *Wss) Wss() string {
	return w.Host + w.Prefix
}

var AppSetting = &App{
	// Env:     "development",
	Env:      "production",
	AppName:  "go-starter",
	AppPort:  "8080",
	TimeZone: map[bool]string{true: os.Getenv("TZ"), false: "UTC"}[os.Getenv("TZ") != ""],
}

var LogSettings = &LogConf{
	Output: "console",
}

var MongoDBSettings = &gmongo.MongoDBConf{
	Enable:   false,
	URI:      "mongodb://localhost:27017",
	Database: "my_default_db",
}

func (l *LogConf) Console() bool {
	return strings.Contains(l.Output, "console")
}

func (l *LogConf) File() bool {
	return strings.Contains(l.Output, "file")
}

func (l *LogConf) Kafka() bool {
	return strings.Contains(l.Output, "kafka")
}

func (l *LogConf) LogFile() string {
	// return filepath.Join(utils.TempDir(), AppSetting.AppName, "mylog.txt")
	return filepath.Join(l.FilePath, AppSetting.AppName, "mylog.txt")
}

func (l *LogConf) LumberjackLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   l.LogFile(),
		MaxSize:    1, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
		Compress:   true,
	}
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
	SubChannels:    "",
}

var Log *logrus.Entry
var filename string = ".env"

// LoadEnvVars will load a ".env[.development|.test]" file if it exists and set ENV vars.
// Useful in development and test modes. Not used in production.
func init() {

	os.Setenv("TZ", AppSetting.TimeZone)

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
	loadLoggerSettings(v1, filename)

	// init log configuration
	InitLog()

	loadWssSettings(v1, filename)
	loadRedisSettings(v1, filename)
	loadMongoDBSettings(v1, filename)
	loadCookieSettings(v1, filename)

}

func InitLog() {

	myLog := logrus.New()

	if LogSettings.File() && LogSettings.Console() {

		lumberjackLogger := LogSettings.LumberjackLogger()

		// assign it to the standard logger
		myLog.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogger))
		// gin log config
		gin.DisableConsoleColor() //  Disable Console Color
		gin.DefaultWriter = io.MultiWriter(lumberjackLogger, os.Stdout)

	} else if LogSettings.Console() {

		// assign it to the standard logger
		myLog.SetOutput(io.MultiWriter(os.Stdout))

		// gin log config
		gin.DisableConsoleColor() //  Disable Console Color
		gin.DefaultWriter = io.MultiWriter(os.Stdout)

	} else if LogSettings.File() {

		lumberjackLogger := LogSettings.LumberjackLogger()

		// assign it to the standard logger
		myLog.SetOutput(io.MultiWriter(lumberjackLogger))
		// gin log config
		gin.DisableConsoleColor() //  Disable Console Color
		gin.DefaultWriter = io.MultiWriter(lumberjackLogger)

	} else {
		myLog.Out = ioutil.Discard
	}

	myLog.SetLevel(logrus.DebugLevel)
	myLog.SetReportCaller(true)

	myLog.SetFormatter(&MyTxtFormatter{
		TimestampFormat: timeStampFormat,
		LogFormat:       "%time% [%lvl%] - %file% > %msg%\n",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			lineMessage := fmt.Sprintf("%s() %s:%d", frame.Function, path.Base(frame.File), frame.Line)
			lineLength := len(lineMessage)
			lineMaxLength := 36
			if lineLength > lineMaxLength {
				lineMessage = "....." + string(lineMessage[lineLength-lineMaxLength+4:lineLength-1])
			} else if lineLength < lineMaxLength {
				lineMessage = utils.PadLeft(lineMessage, " ", lineMaxLength)
			}

			return "", "{" + lineMessage + "}"
		},
	})

	if LogSettings.Kafka() {

		kafkaLogTopic := LogSettings.KafkaTopic
		kafkaServerAddr := LogSettings.KafkaServer

		hook, err := lkh.NewKafkaHook(
			"kh",
			[]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel},
			// &logrus.JSONFormatter{},
			&MyJSONFormatter{
				TimestampFormat: timeStampFormat,
				PrettyPrint:     false,
				AppName:         AppSetting.AppName,
				Env:             AppSetting.Env,
				CapationGen:     1,
				FieldMap: FieldMap{
					"time": "@timestamp",
					"msg":  "@message",
				},
			},
			strings.Split(kafkaServerAddr, ","),
		)

		if err != nil {
			myLog.WithError(err).Error(fmt.Sprintf("Kafka Append Initialization failed: kafkaServer=%s, errorMessage=%s", kafkaServerAddr, err.Error()))
			Log = logrus.NewEntry(myLog)
		} else {
			myLog.Hooks.Add(hook)
			Log = myLog.WithField("topics", []string{kafkaLogTopic})
			Log.Info("Initialization logger completed: kafkaSever=", LogSettings.KafkaServer, ", logTopic=", LogSettings.KafkaTopic)
		}

	} else {
		Log = logrus.NewEntry(myLog)
	}

}

func loadAppSettings(v1 *viper.Viper, filename string) {
	err := v1.Unmarshal(&AppSetting)
	if err != nil {
		log.Println("viper parse AppSettings error: configFile=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		log.Println("AppSetting: TimeZone=" + AppSetting.TimeZone + " Env=" + AppSetting.Env + ", AppName=" + AppSetting.AppName + ", AppPort=" + AppSetting.AppPort + ", appRoot=" + utils.RootDir())
	}
}

func loadLoggerSettings(v1 *viper.Viper, filename string) {
	err := v1.Unmarshal(&LogSettings)
	if err != nil {
		log.Println("viper parse LogSettings error: configFile=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		log.Println("LogSettings: Output=" + LogSettings.Output + ", logFile=" + LogSettings.LogFile())
	}
}

func loadRedisSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&RedisConf)

	if err != nil {
		Log.Info("viper parse RedisConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		if RedisConf.Mode == gredis.Disabled {
			Log.Info("Redis-Not-Enabled: Mode=" + utils.ToString(RedisConf.Mode) + ", Host=" + RedisConf.Host)
		} else {
			gredis.Init(RedisConf, Log)
		}
	}

}
func loadMongoDBSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&MongoDBSettings)

	if err != nil {
		Log.Info("viper parse MongoDBSettings error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		if !MongoDBSettings.Enable {
			Log.Info("MongoDb-Not-Enabled: Enabled=" + utils.ToString(MongoDBSettings.Enable) + ", Host=" + RedisConf.Host)
		} else {
			gmongo.Init(MongoDBSettings, Log)
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
