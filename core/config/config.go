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

	"github.com/chunhui2001/go-starter/core/ges"
	"github.com/chunhui2001/go-starter/core/ghttp"
	"github.com/chunhui2001/go-starter/core/gmongo"
	"github.com/chunhui2001/go-starter/core/gredis"
	"github.com/chunhui2001/go-starter/core/gsql"
	_ "github.com/chunhui2001/go-starter/core/gzk"
	"github.com/chunhui2001/go-starter/core/utils"
	_ "github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"

	lkh "github.com/gfremex/logrus-kafka-hook"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var timeStampFormat = "2006-01-02T15:04:05.000Z07:00"
var configLoggerLines [][]string = [][]string{}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().Format(timeStampFormat) + " [STDOUT] - " + string(bytes))
}

type AppConf struct {
	Env        string `mapstructure:"GIN_ENV"`
	AppName    string `mapstructure:"APP_NAME"`
	AppPort    string `mapstructure:"APP_PORT"`
	TimeZone   string `mapstructure:"APP_TIMEZONE"`
	DemoEnable bool   `mapstructure:"ENABLE_DEMO"`
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
	MaxAge int    `mapstructure:"COOKIE_MAXAGE"`
}

type LogConf struct {
	Output      string `mapstructure:"LOG_OUTPUT"`
	FilePath    string `mapstructure:"LOG_FILE_PATH"`
	KafkaServer string `mapstructure:"LOG_KAFKA_SERVER"`
	KafkaTopic  string `mapstructure:"LOG_KAFKA_TOPIC"`
}

type WebPageConf struct {
	Enable    bool   `mapstructure:"WEB_PAGE_ENABLE"`
	Root      string `mapstructure:"WEB_PAGE_ROOT"`
	Master    string `mapstructure:"WEB_PAGE_MASTER"`
	Extension string `mapstructure:"WEB_PAGE_Extension"`
	LoginUrl  string `mapstructure:"WEB_PAGE_LOGIN"`
	SignUpUrl string `mapstructure:"WEB_PAGE_SIGNUP"`
}

func (w *Wss) Wss() string {
	return w.Host + w.Prefix
}

var AppSetting = &AppConf{
	// Env:     "development",
	Env:        "production",
	AppName:    "go-starter",
	AppPort:    "8080",
	TimeZone:   map[bool]string{true: os.Getenv("TZ"), false: "UTC"}[os.Getenv("TZ") != ""],
	DemoEnable: true,
}

var LogSettings = &LogConf{
	Output: "console",
}

var WebPageSettings = &WebPageConf{
	Enable:    false,
	Root:      "views",
	Master:    "layouts/master",
	Extension: ".html",
	LoginUrl:  "/login",
	SignUpUrl: "/signup",
}

var MongoDBSettings = &gmongo.MongoDBConf{
	Enable:   false,
	URI:      "mongodb://localhost:27017",
	Database: "my_default_db",
}

var EsSettings = &ges.ESConf{
	Enable:  false,
	Servers: "http://localhost:9200",
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
	return filepath.Join(l.FilePath, AppSetting.AppName, "applog.txt")
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

var MySqlConf = &gsql.MySql{
	Enable: false,
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
	var envfile = ".env/.env." + env

	var WEB_PAGE_Enable string = os.Getenv("WEB_PAGE_Enable")
	log.Println(WEB_PAGE_Enable)

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

	loadWebPageSettings(v1, filename)
	loadWssSettings(v1, filename)
	loadRedisSettings(v1, filename)
	loadMongoDBSettings(v1, filename)
	loadCookieSettings(v1, filename)
	loadMySqlSettings(v1, filename)
	loadEsSettings(v1, filename)

	printConfigLogLines()

	ghttp.Init(Log)

}

func printConfigLogLines() {

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	for _, v := range configLoggerLines {
		table.Append(v)
	}

	table.Render() // Send output
	tableLines := utils.Split(tableString.String(), "\n")

	for _, line := range tableLines {
		if line != "" {
			Log.Info(line)
		}
	}
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
			lineMessage := fmt.Sprintf("%s() %s(%d):%d", frame.Function, path.Base(frame.File), utils.GoroutineId(), frame.Line)
			lineLength := len(lineMessage)
			lineMaxLength := 36
			if lineLength > lineMaxLength {
				lineMessage = "....." + string(lineMessage[lineLength-lineMaxLength+4:lineLength])
			} else if lineLength < lineMaxLength {
				lineMessage = utils.PadLeft(lineMessage, " ", lineMaxLength+1)
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
		log.Println("AppSetting: TimeZone=" + AppSetting.TimeZone + ", GIN_ENV=" + AppSetting.Env + ", AppName=" + AppSetting.AppName + ", AppPort=" + AppSetting.AppPort + ", appRoot=" + utils.RootDir())
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
		configLoggerLines = append(configLoggerLines, []string{"RedisSettings", "Mode=" + RedisConf.Mode.String()})
		if !RedisConf.Disabled() {
			gredis.Init(RedisConf, Log)
		}
	}

}

func loadMySqlSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&MySqlConf)

	if err != nil {
		Log.Info("viper parse MySqlConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		configLoggerLines = append(configLoggerLines, []string{"MySqlSettings", "Enabled=" + utils.ToString(MySqlConf.Enable)})
		if MySqlConf.Enable != false {
			gsql.Init(MySqlConf, Log)
		}
	}

}

func loadEsSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&EsSettings)

	if err != nil {
		Log.Info("viper parse ESConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		configLoggerLines = append(configLoggerLines, []string{"EsSettings", "Enabled=" + utils.ToString(EsSettings.Enable)})
		if EsSettings.Enable != false {
			ges.Init(EsSettings, Log)
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
		configLoggerLines = append(configLoggerLines, []string{"MongoDbSettings", "Enabled=" + utils.ToString(MongoDBSettings.Enable)})
		if MongoDBSettings.Enable {
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
		if WssSetting.Enable {
			configLoggerLines = append(configLoggerLines, []string{"WssSetting", "Enable=" + utils.ToString(WssSetting.Enable) + ", Host=" + WssSetting.Host + ", Prefix=" + WssSetting.Prefix})
		} else {
			configLoggerLines = append(configLoggerLines, []string{"WssSetting", "Enable=" + utils.ToString(WssSetting.Enable)})
		}
	}

}

func loadWebPageSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&WebPageSettings)

	if err != nil {
		Log.Info("viper parse WebPageSettings error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		if WebPageSettings.Enable {
			configLoggerLines = append(configLoggerLines, []string{"WebPageSettings", "Enable=" + utils.ToString(WebPageSettings.Enable) + ", Root=" + filepath.Join(utils.RootDir(), WebPageSettings.Root)})
		} else {
			configLoggerLines = append(configLoggerLines, []string{"WebPageSettings", "Enable=" + utils.ToString(WebPageSettings.Enable)})
		}
	}

}

func loadCookieSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&CookieSetting)

	if err != nil {
		Log.Info("viper parse CookieSetting error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		if CookieSetting.Enable {
			configLoggerLines = append(configLoggerLines, []string{"CookieSetting", "Enabled=" + strconv.FormatBool(CookieSetting.Enable) + ", Name=" + CookieSetting.Name + ", Secret=" + CookieSetting.Secret + ", MaxAge=" + fmt.Sprint(CookieSetting.MaxAge)})
		} else {
			configLoggerLines = append(configLoggerLines, []string{"CookieSetting", "Enabled=" + strconv.FormatBool(CookieSetting.Enable)})
		}

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
	v.AutomaticEnv() // 将读取当前目录下的 .env 配置文件或"环境变量", .env 优先级最高
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
