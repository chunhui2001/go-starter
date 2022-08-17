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
var RedisConf = &Redis{
	Enable: false,
}

var filename string = ".env"

// LoadEnvVars will load a ".env[.development|.test]" file if it exists and set ENV vars.
// Useful in development and test modes. Not used in production.
func init() {

	var env string = os.Getenv("GIN_ENV")
	var envfile = ".env." + env

	if exists, _ := utils.FileExists(utils.RootDir() + "/" + envfile); exists == true {
		filename = envfile
	} else {
		log.Println("Configuration loading " + utils.RootDir() + "/" + envfile + " file error, use .env file.")
	}

	err := godotenv.Load(filename)

	if err != nil {
		log.Println("Configuration loading " + utils.RootDir() + "/" + filename + " file error, errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	}

	v1 := readConfig(filename, map[string]interface{}{})

	loadAppSettings(v1, filename)
	loadRedisSettings(v1, filename)
	loadWssSettings(v1, filename)

	log.Println("Configuration loaded " + utils.RootDir() + "/" + filename + " successful.")

}

func loadAppSettings(v1 *viper.Viper, filename string) {
	err := v1.Unmarshal(&AppSetting)
	if err != nil {
		log.Println("viper parse AppSettings error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	}
}

func loadRedisSettings(v1 *viper.Viper, filename string) {

	REDIS_Enable, _ := strconv.ParseBool(GetEnv("REDIS_Enable", "false"))

	if REDIS_Enable {
		err := v1.Unmarshal(&RedisConf)
		if err != nil {
			log.Println("viper parse RedisConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
			os.Exit(3)
			return
		}
	} else {
		RedisConf = &Redis{Enable: false}
	}

}

func loadWssSettings(v1 *viper.Viper, filename string) {

	WSS_PREFIX := GetEnv("WSS_PREFIX", "")

	if WSS_PREFIX != "" {
		err := v1.Unmarshal(&WssSetting)
		if err != nil {
			log.Println("viper parse WssSetting error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
			os.Exit(3)
			return
		} else {
			log.Println("WssSetting.Prefix: file=" + WssSetting.Prefix)
		}
	} else {
		WssSetting = &Wss{Prefix: ""}
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
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		log.Println("viper loaded error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return nil
	}
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
