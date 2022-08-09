package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chunhui2001/go-starter/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	Env     string `mapstructure:"GIN_ENV"`
	AppName string `mapstructure:"APP_NAME"`
	AppPort string `mapstructure:"APP_PORT"`
}

type Wss struct {
	Prefix string `mapstructure:"WSS_PREFIX"`
}

type Redis struct {
	Enable      bool          `mapstructure:"REDIS_Enable"`
	Host        string        `mapstructure:"REDIS_Host"`
	Passwd      string        `mapstructure:"REDIS_Password"`
	MaxIdle     int           `mapstructure:"REDIS_MaxIdle"`
	MaxActive   int           `mapstructure:"REDIS_MaxActive"`
	IdleTimeout time.Duration `mapstructure:"REDIS_IdleTimeout"`
}

var AppSetting = &App{}

var WssSetting = &Wss{}

var RedisConf = &Redis{}

var filename string = ".env"

// LoadEnvVars will load a ".env[.development|.test]" file if it exists and set ENV vars.
// Useful in development and test modes. Not used in production.
func init() {

	var env string = os.Getenv("GIN_ENV")
	var envfile = utils.RootDir() + "/" + ".env." + env

	if exists, _ := utils.Exists(envfile); exists == true {
		filename = envfile
	} else {
		log.Println("loading " + envfile + " file error, use .env file.")
	}

	err := godotenv.Load(filename)

	if err != nil {
		log.Println("loading " + filename + " file error, errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	}

	viper.AddConfigPath(utils.RootDir())
	viper.SetConfigName(filename)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		log.Println("viper loaded error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	}

	loadAppSettings(filename)
	loadRedisSettings(filename)
	loadWssSettings(filename)

	log.Println("viper " + filename + " loaded.")

}

func loadAppSettings(filename string) {
	err := viper.Unmarshal(&AppSetting)
	if err != nil {
		log.Println("viper parse AppSettings error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		log.Println("viper readed AppSettings appName=" + AppSetting.AppName)
	}
}

func loadRedisSettings(filename string) {

	REDIS_Enable, _ := strconv.ParseBool(GetEnv("REDIS_Enable", "false"))

	if REDIS_Enable {
		err := viper.Unmarshal(&RedisConf)
		if err != nil {
			log.Println("viper parse RedisConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
			os.Exit(3)
			return
		} else {
			log.Println("viper readed RedisConf redisHost=" + RedisConf.Host)
		}
	} else {
		RedisConf = &Redis{Enable: false}
	}

}

func loadWssSettings(filename string) {

	WSS_PREFIX := GetEnv("WSS_PREFIX", "")

	if WSS_PREFIX != "" {
		err := viper.Unmarshal(&WssSetting)
		if err != nil {
			log.Println("viper parse WssSetting error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
			os.Exit(3)
			return
		} else {
			log.Println("viper readed WssSetting wssPrefix=" + WssSetting.Prefix)
		}
	} else {
		WssSetting = &Wss{Prefix: ""}
	}

}

// GetEnv returns an environment variable or a default value if not present
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
