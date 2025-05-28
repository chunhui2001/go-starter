package config

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"io"
	"path"
	"path/filepath"
	"runtime"

	"github.com/chunhui2001/go-starter/core/built"
	"github.com/chunhui2001/go-starter/core/ges"
	"github.com/chunhui2001/go-starter/core/ghttp"
	"github.com/chunhui2001/go-starter/core/gid"
	"github.com/chunhui2001/go-starter/core/gmongo"
	"github.com/chunhui2001/go-starter/core/goes"
	"github.com/chunhui2001/go-starter/core/googleapi"
	"github.com/chunhui2001/go-starter/core/grabbit"
	"github.com/chunhui2001/go-starter/core/gredis"
	"github.com/chunhui2001/go-starter/core/grtask"
	"github.com/chunhui2001/go-starter/core/gsql"
	"github.com/chunhui2001/go-starter/core/gzok"
	"github.com/chunhui2001/go-starter/core/gztask"
	"github.com/chunhui2001/go-starter/core/utils"
	_ "github.com/joho/godotenv"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/gin-gonic/gin"

	lkh "github.com/gfremex/logrus-kafka-hook"
	"github.com/jinzhu/copier"
	// "github.com/olekukonko/tablewriter"
	"github.com/rifflock/lfshook"
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
	NodeId     int64  `mapstructure:"NODE_ID"`
	Env        string `mapstructure:"GIN_ENV"`
	AppName    string `mapstructure:"APP_NAME"`
	AppPort    string `mapstructure:"APP_PORT"`
	TimeZone   string `mapstructure:"APP_TIMEZONE"`
	DemoEnable bool   `mapstructure:"ENABLE_DEMO"`
	AppVersion string `mapstructure:"APP_VERSION"`
	OS         string `mapstructure:"APP_OS"`
	CaptainGEN int    `mapstructure:"CAPTAIN_GEN"`
}

type Wss struct {
	Enable       bool   `mapstructure:"WSS_ENABLE"`
	Prefix       string `mapstructure:"WSS_PREFIX"`
	Host         string `mapstructure:"WSS_HOST"`
	PrintMessage bool   `mapstructure:"WSS_PRINT_MESSAGE"`
}

type Cookie struct {
	Enable bool   `mapstructure:"COOKIE_ENABLE"`
	Name   string `mapstructure:"COOKIE_NAME"`
	Secret string `mapstructure:"COOKIE_SECRET"`
	MaxAge int    `mapstructure:"COOKIE_MAXAGE"`
}

type LogConf struct {
	Output           string `mapstructure:"LOG_OUTPUT"`
	FilePath         string `mapstructure:"LOG_FILE_PATH"`
	FileMaxSize      int    `mapstructure:"LOG_FILE_MAX_SIZE"`
	FileMaxBackups   int    `mapstructure:"LOG_FILE_MAX_BACKUPS"`
	FileMaxAge       int    `mapstructure:"LOG_FILE_MAX_AGE"`
	KafkaServer      string `mapstructure:"LOG_KAFKA_SERVER"`
	KafkaTopic       string `mapstructure:"LOG_KAFKA_TOPIC"`
	FileFormatter    string `mapstructure:"LOG_FILE_FORMATTER"`    // json OR txt
	ConsoleFormatter string `mapstructure:"LOG_CONSOLE_FORMATTER"` // json OR txt
}

type WebPageConf struct {
	Enable    bool   `mapstructure:"WEB_PAGE_ENABLE"`
	Root      string `mapstructure:"WEB_PAGE_ROOT"`
	Master    string `mapstructure:"WEB_PAGE_MASTER"`
	Extension string `mapstructure:"WEB_PAGE_Extension"`
	LoginUrl  string `mapstructure:"WEB_PAGE_LOGIN"`
	SignUpUrl string `mapstructure:"WEB_PAGE_SIGNUP"`
}

type GraphServerConf struct {
	Enable        bool   `mapstructure:"GRAPHQL_ENABLE"`
	ServerURi     string `mapstructure:"GRAPHQL_SERVER_URI"`
	PlayGroundURi string `mapstructure:"GRAPHQL_PLAYGROUND_URI"`
}

func (w *Wss) Wss() string {
	if w.Enable {
		return w.Host + w.Prefix
	}
	return ""
}

var AppSetting = &AppConf{
	Env:        "production",
	AppName:    "go-starter",
	AppPort:    "0.0.0.0:8080",
	TimeZone:   map[bool]string{true: os.Getenv("TZ"), false: "UTC"}[os.Getenv("TZ") != ""],
	NodeId:     utils.Ip2Int(utils.OutboundIP()) % 1023,
	DemoEnable: false,
}

var GraphServerSetting = &GraphServerConf{
	Enable: false,
}

var SimpleGTaskConf = &gztask.SimpleGTask{
	Enable:   false,
	ID:       "g4qUY4f17Bk",
	Name:     "一个示例定时任务执行",
	Expr:     "* * * * * *",
	PrintLog: true,
}

var LogSettings = &LogConf{
	Output:           "console",
	FileMaxSize:      1, // 1MB
	FileMaxBackups:   10,
	FileMaxAge:       30,
	FileFormatter:    "txt", // json OR txt,
	ConsoleFormatter: "txt", // json OR txt
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
	Enable:      false,
	Servers:     "http://localhost:9200",
	PrettyPrint: false,
}

var OpenEsSettings = &goes.OpenESConf{
	Enable:      false,
	Servers:     "http://localhost:9200",
	PrettyPrint: false,
}

var GZokConf = &gzok.GZokConf{
	Enabled: false,
	Hosts:   []string{"127.0.0.1:2181"},
}

var jsonFormatter = func() *MyJSONFormatter {
	return &MyJSONFormatter{
		TimestampFormat: timeStampFormat,
		PrettyPrint:     false,
		AppName:         AppSetting.AppName,
		Env:             AppSetting.Env,
		CaptainGEN:      AppSetting.CaptainGEN,
		IP:              utils.OutboundIP().String(),
		FieldMap: FieldMap{
			"time": "@timestamp",
			"msg":  "message",
		},
	}
}

var txtFormatter = func() *MyTxtFormatter {
	return &MyTxtFormatter{
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
	}
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
		MaxSize:    l.FileMaxSize, // MB
		MaxBackups: l.FileMaxBackups,
		MaxAge:     l.FileMaxAge, // days
		Compress:   true,
	}
}

var WssSetting = &Wss{
	Enable:       false,
	Prefix:       "",
	PrintMessage: true,
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
	PrintMessage:   true,
}

var MySqlConf = &gsql.MySql{
	Enable: false,
}

var RabbitMQConf = &grabbit.GRabbitConf{
	Enable: false,
}

var HttpClientConf = &ghttp.HttpConf{
	Timeout:             150, // * time.Second
	IdleConnTimeout:     90,
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 100,
	MaxConnsPerHost:     100,
	PrintCurl:           true,
	PrintDebug:          false,
}

var GoogleAPIConfSettings = &googleapi.GoogleAPIConf{
	Enable:          false,
	CredentialsFile: "resources/googleapi-oauth-credentials.json",
	TokenFile:       "resources/googleapi-oauth-token.json",
	Scopes:          []string{"https://www.googleapis.com/auth/drive", "https://www.googleapis.com/auth/drive.file", "https://www.googleapis.com/auth/drive.metadata", "https://www.googleapis.com/auth/drive.appdata", "https://www.googleapis.com/auth/spreadsheets"},
}

var Log *logrus.Entry
var myViper *viper.Viper
var filename string = ".env"
var applicationConfig map[string]interface{}

func AppRoot() string {
	return utils.RootDir()
}

// GetEnv returns an environment variable or a default value if not present
func GetEnv(key, defaultValue string) string {

	value := os.Getenv(key)

	if value != "" {
		return value
	}

	value = myViper.GetString(key)

	if value != "" {
		return value
	}

	return defaultValue

}

// LoadEnvVars will load a ".env[.development|.test]" file if it exists and set ENV vars.
// Useful in development and test modes. Not used in production.
func init() {

	cmdArgs := os.Args

	if len(cmdArgs) > 1 {
		if cmdArgs[1] == "version" {
			fmt.Println(built.INFO.Info())
			os.Exit(0)
		}
	}

	os.Setenv("TZ", AppSetting.TimeZone)

	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	var env string = os.Getenv("GIN_ENV")
	var defaultenv = ".env/.env"
	var envfile = ".env/.env." + env

	if exists, _ := utils.FileExists(filepath.Join(AppRoot(), envfile)); exists {
		filename = envfile
	} else {
		if env == "" {
			filename = ""
		} else {
			log.Println("Configuration loading " + filepath.Join(AppRoot(), envfile) + " file error")
		}
	}

	log.Println(built.INFO.Info())

	if v1 := readConfig(map[string]interface{}{}, defaultenv, filename); v1 != nil {

		myViper = v1

		loadAppSettings(v1, filename)
		loadLoggerSettings(v1, filename)

		// init log configuration
		InitLog()

		loadWebPageSettings(v1, filename)
		loadWssSettings(v1, filename)
		loadHttpClientSettings(v1, filename)
		loadRedisSettings(v1, filename)
		loadZookeeperSettings(v1, filename)
		loadMongoDBSettings(v1, filename)
		loadCookieSettings(v1, filename)
		loadMySqlSettings(v1, filename)
		loadEsSettings(v1, filename)
		loadOpenEsSettings(v1, filename)
		loadSimpleGTaskSettings(v1, filename)
		loadRabbitSettings(v1, filename)
		loadGraphServerSettings(v1, filename)
		loadGoogleApiServiceSettings(v1, filename)

		printConfigLogLines()

		grtask.Init(Log, AppSetting.NodeId)
		gztask.Init(Log, AppSetting.AppName, SimpleGTaskConf)

		gid.Init(Log, AppSetting.NodeId)

		loadYamlConfiguraion()

	} else {
		Log = logrus.NewEntry(logrus.New())
	}

}

func printConfigLogLines() {

	// tableString := &strings.Builder{}
	// table := tablewriter.NewWriter(tableString)
	// table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	// table.SetRowLine(true)
	// table.SetAutoWrapText(false)
	// table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	// for _, v := range configLoggerLines {
	// 	table.Append(v)
	// }

	// table.Render() // Send output
	// tableLines := utils.Split(tableString.String(), "\n")

	// for _, line := range tableLines {
	// 	if line != "" {
	// 		Log.Info(line)
	// 	}
	// }
}

func InitLog() {

	myLog := logrus.New()

	if LogSettings.Console() {

		// assign it to the standard logger
		myLog.SetOutput(io.MultiWriter(os.Stdout))

		// gin log config
		gin.DisableConsoleColor() //  Disable Console Color
		gin.DefaultWriter = io.MultiWriter(os.Stdout)

		myLog.SetLevel(logrus.DebugLevel)
		myLog.SetReportCaller(true)

		if LogSettings.ConsoleFormatter == "json" {
			myLog.SetFormatter(jsonFormatter())
		} else {
			myLog.SetFormatter(txtFormatter())
		}

	} else {
		myLog.Out = ioutil.Discard
	}

	if LogSettings.File() {

		writer := LogSettings.LumberjackLogger()

		if LogSettings.FileFormatter == "json" {
			myLog.Hooks.Add(lfshook.NewHook(
				lfshook.WriterMap{logrus.InfoLevel: writer, logrus.ErrorLevel: writer},
				jsonFormatter(),
			))
		} else {
			myLog.Hooks.Add(lfshook.NewHook(
				lfshook.WriterMap{logrus.InfoLevel: writer, logrus.ErrorLevel: writer},
				txtFormatter(),
			))
		}

	}

	if LogSettings.Kafka() {

		kafkaLogTopic := LogSettings.KafkaTopic
		kafkaServerAddr := LogSettings.KafkaServer

		hook, err := lkh.NewKafkaHook(
			"kh",
			[]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.DebugLevel},
			jsonFormatter(), // &logrus.JSONFormatter{},
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
		AppSetting.AppVersion = built.INFO.Commit
		AppSetting.OS = built.INFO.OS
		log.Println("viper parse AppSettings error: Version=" + AppSetting.AppVersion + ", OS=" + AppSetting.OS + ", configFile=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {

		v, _ := mem.VirtualMemory()

		infoStat, _ := cpu.Info()
		cpuMode := infoStat[len(infoStat)-1].ModelName
		physicalID := infoStat[len(infoStat)-1].PhysicalID
		cores := infoStat[len(infoStat)-1].Cores

		log.Println("InfoStat:" +
			" Hostname=" + utils.Hostname() +
			", OutboundIP=" + utils.OutboundIP().String() +
			", NumCPUs=" + utils.ToString(runtime.NumCPU()) +
			", Cores=" + utils.ToString(cores) +
			", TotalMem=" + utils.HumanFileSizeUint(v.Total) +
			", PhysicalID=" + physicalID +
			", GOMAXPROCS=" + os.Getenv("GOMAXPROCS") +
			", CPUMode=" + cpuMode)

		AppSetting.AppVersion = built.INFO.Commit
		AppSetting.OS = built.INFO.OS

		log.Println("AppSetting:" +
			" TimeZone=" + AppSetting.TimeZone +
			", GIN_ENV=" + AppSetting.Env +
			", NODE_ID=" + utils.ToString(AppSetting.NodeId) +
			", AppName=" + AppSetting.AppName +
			", AppPort=" + AppSetting.AppPort +
			", Version=" + AppSetting.AppVersion +
			", OS=" + AppSetting.OS +
			", AppRoot=" + AppRoot())
	}

}

func loadLoggerSettings(v1 *viper.Viper, filename string) {
	err := v1.Unmarshal(&LogSettings)
	if err != nil {
		log.Println("viper parse LogSettings error: configFile=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		log.Println("LogSettings: Output=" + LogSettings.Output + ", Formatter=" + LogSettings.FileFormatter + ", LogFile=" + LogSettings.LogFile() + ", MaxSize=" + utils.ToString(LogSettings.FileMaxSize) + "mb")
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

func loadZookeeperSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&GZokConf)

	if err != nil {
		Log.Info("viper parse GZokConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		configLoggerLines = append(configLoggerLines, []string{"Zookeeper", "Enable=" + fmt.Sprint(GZokConf.Enabled)})
		if GZokConf.Enabled {
			gzok.Init(GZokConf, Log)
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
		if MySqlConf.Enable {
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
		if EsSettings.Enable {
			ges.Init(EsSettings, Log)
		}
	}

}

func loadOpenEsSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&OpenEsSettings)

	if err != nil {
		Log.Info("viper parse OpenEsSettings error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		configLoggerLines = append(configLoggerLines, []string{"OpenEsSettings", "Enabled=" + utils.ToString(OpenEsSettings.Enable) + ", PrettyPrint=" + utils.ToString(OpenEsSettings.PrettyPrint)})
		if OpenEsSettings.Enable {
			goes.Init(OpenEsSettings, Log)
		}
	}

}

func loadSimpleGTaskSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&SimpleGTaskConf)

	if err != nil {
		Log.Info("viper parse SimpleGTaskConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		configLoggerLines = append(configLoggerLines, []string{"SimpleGTaskConf", "Enabled=" + utils.ToString(SimpleGTaskConf.Enable)})
	}

}

func loadRabbitSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&RabbitMQConf)

	if err != nil {
		Log.Info("viper parse SimpleGTaskConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {

		configLoggerLines = append(configLoggerLines, []string{"RabbitMQConf", "Enabled=" + utils.ToString(RabbitMQConf.Enable)})

		if RabbitMQConf.Enable {
			grabbit.Init(RabbitMQConf, Log)
		}

	}

}

func loadGraphServerSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&GraphServerSetting)

	if err != nil {
		Log.Info("viper parse GraphServerSetting error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		configLoggerLines = append(configLoggerLines, []string{"GraphServer", "Enabled=" + utils.ToString(GraphServerSetting.Enable)})
	}

}

func loadGoogleApiServiceSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&GoogleAPIConfSettings)

	if err != nil {
		Log.Info("viper parse GoogleAPIConfSettings  error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		configLoggerLines = append(configLoggerLines, []string{"GoogleAPIConfSettings", "Enabled=" + utils.ToString(GoogleAPIConfSettings.Enable) + ",CREDENTIALS_FILE=" + GoogleAPIConfSettings.CredentialsFile})
		if GoogleAPIConfSettings.Enable {
			googleapi.Init(GoogleAPIConfSettings, Log)
		}
	}

}

func loadHttpClientSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&HttpClientConf)

	if err != nil {
		Log.Info("viper parse HttpClientConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		return
	} else {
		ghttp.Init(HttpClientConf, Log)
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
			configLoggerLines = append(configLoggerLines, []string{"WebPageSettings", "Enable=" + utils.ToString(WebPageSettings.Enable) + ", Root=" + filepath.Join(AppRoot(), WebPageSettings.Root)})
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

func loadYamlConfiguraion() {
	var agolloService string = os.Getenv("APOLLO_CONFIGSERVICE")

	if len(agolloService) > 0 {
		readApollo2()

		return
	}

	var env string = os.Getenv("GIN_ENV")

	var f = func(file string) map[string]interface{} {
		yamlFilePath := filepath.Join(utils.RootDir(), ".env", file)

		if exists, _ := utils.FileExists(yamlFilePath); exists {
			if yamlContent, err := utils.ReadFile(yamlFilePath); err != nil {
				Log.Errorf(`Read-Yaml-File-Error: FilePath=%s, ErrorMessage=%s`, yamlFilePath, err.Error())
			} else {
				var body map[string]interface{}
				if err := yaml.Unmarshal([]byte(yamlContent), &body); err != nil {
					Log.Errorf(`Loading-Yaml-File-Error: FilePath=%s, ErrorMessage=%s`, yamlFilePath, err.Error())
				} else {
					Log.Infof(`LoadedYaml-Configuration: FilePath=%s`, yamlFilePath)
					return body
				}
			}
		}

		return nil
	}

	var body1 map[string]interface{} = f("application.yml")
	var body2 map[string]interface{} = f("application-" + env + ".yml")

	if body1 != nil {
		if err := copier.CopyWithOption(&applicationConfig, &body1, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
			panic(err)
		}
	}

	if body2 != nil {
		if err := copier.CopyWithOption(&applicationConfig, &body2, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
			panic(err)
		}
	}
}

// yaml config
func ReadConfig(key string, data any) error {
	value := applicationConfig[key]

	if err := json.Unmarshal(utils.ToJsonBytes(value), data); err != nil {
		Log.Errorf("Configuration-Error: Key=%s, ErrorMessage=%v", "key", err)
		return err
	}

	return nil
}

func AllConfig() map[string]interface{} {

	return applicationConfig
}

func readConfig(defaults map[string]interface{}, filenames ...string) *viper.Viper {
	v := viper.New()

	var agolloService string = os.Getenv("APOLLO_CONFIGSERVICE")

	if len(agolloService) > 0 {
		readApollo1(v)

		return v
	}

	var f = func(file string, defaultMaps map[string]interface{}) *viper.Viper {
		v := viper.New()

		for key, value := range defaultMaps {
			v.SetDefault(key, value)
		}

		v.SetConfigName(file)
		v.SetConfigType("env")
		// v.AddConfigPath("/etc/appname/")   // path to look for the config file in
		// v.AddConfigPath("$(home)/.env") // call multiple times to add many search paths
		v.AddConfigPath(AppRoot())
		v.AutomaticEnv() // 将读取当前目录下的 .env 配置文件或"环境变量", .env 优先级最高

		err := v.ReadInConfig()

		if err != nil {
			log.Println("viper loaded error: AppRoot=" + AppRoot() + ", file=" + file + ", errorMessage=" + fmt.Sprint(err) + ".")
			return nil
		}

		log.Println("viper Configuration loaded " + filepath.Join(AppRoot(), file) + " successful.")

		return v
	}

	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	for _, fname := range filenames {
		if fname != "" {
			v = f(fname, v.AllSettings())
		}
	}

	return v
}

func readApollo1(v *viper.Viper) {

	var agolloService string = os.Getenv("APOLLO_CONFIGSERVICE")
	var appId string = os.Getenv("APP_ID")
	var env string = os.Getenv("GIN_ENV") // cluster

	if len(agolloService) <= 0 {
		log.Println("Loaded-agollo-properties processed: Disabled")
		return
	}

	var s1 = fmt.Sprintf("%s/configs/%s/%s/application.properties", agolloService, appId, env)

	httpResponse1 := sendApolloRequest(s1)

	if httpResponse1 != nil {
		responseMap := utils.AsMap(httpResponse1)
		config := responseMap["configurations"].(map[string]interface{})

		for key, val := range config {
			v.SetDefault(strings.TrimSpace(key), strings.TrimSpace(val.(string)))
		}

		log.Printf("Loaded-agollo-properties Successful: configKeyLen=%d, agolloService=%s", len(config), s1)
	}
}

func readApollo2() {

	var agolloService string = os.Getenv("APOLLO_CONFIGSERVICE")
	var appId string = os.Getenv("APP_ID")
	var env string = os.Getenv("GIN_ENV") // cluster

	if len(agolloService) <= 0 {
		log.Println("Loaded-agollo-properties processed: Disabled")
		return
	}

	var s2 = fmt.Sprintf("%s/configs/%s/%s/application.yaml", agolloService, appId, env)

	httpResponse2 := sendApolloRequest(s2)

	if httpResponse2 != nil {
		responseMap := utils.AsMap(httpResponse2)
		config := responseMap["configurations"].(map[string]interface{})["content"].(string)

		var body map[string]interface{}

		if err := yaml.Unmarshal([]byte(config), &body); err != nil {
			Log.Errorf(`Loading-Yaml-File-Error: s2=%s, ErrorMessage=%s`, s2, err.Error())
		}

		if body != nil {
			if err := copier.CopyWithOption(&applicationConfig, &body, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
				panic(err)
			}
		}

		Log.Infof("Loaded-agollo-yaml Successful: configKeyLen=%d, agolloService=%s", len(body), s2)
	}
}

func sendApolloRequest(url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)

	res, err2 := (&http.Client{}).Do(req)

	if err2 != nil {
		log.Printf("Loaded-agollo-properties Error: Url=%s, ErrorMessage=%s", url, err2)

		return nil
	}

	var headerKey string = os.Getenv("APOLLO_HEADER_KEY")
	var secretKey string = os.Getenv("APOLLO_SECRET_KEY")

	log.Printf("sendApolloRequest: headerKey=%s, exists=%s", headerKey, res.Header.Get(headerKey))

	if res.Header.Get(headerKey) == "true" {

		log.Printf("sendApolloRequest: secretKeyLength=%d", len(secretKey))

		// 解密
		raw, errr1 := base64.RawStdEncoding.DecodeString(secretKey)

		if errr1 != nil {
			log.Printf("sendApolloRequest-DecodeString: Error: ErrorMessage=%s", errr1)
		}

		sk, errr2 := x509.ParsePKCS8PrivateKey(raw)

		if errr2 != nil {
			log.Printf("sendApolloRequest-ParsePKCS8PrivateKey: Error: ErrorMessage=%s", errr2)
		}

		privKey := sk.(*rsa.PrivateKey)
		partLen := privKey.PublicKey.N.BitLen() / 8

		resBody, errr3 := io.ReadAll(res.Body)

		if errr3 != nil {
			log.Printf("sendApolloRequest-ParsePKCS8PrivateKey: Error: ErrorMessage=%s", errr3)
		}

		chunks := split(resBody, partLen)

		buffer := bytes.NewBufferString("")

		for _, chunk := range chunks {
			decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, privKey, chunk)

			buffer.Write(decrypted)
		}

		return buffer.Bytes()
	}

	resBody, _ := io.ReadAll(res.Body)

	return resBody
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)

	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}

	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}

	return chunks
}
