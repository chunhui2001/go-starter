package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chunhui2001/go-starter/utils"
	_ "github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	Env     string `mapstructure:"GIN_ENV"`
	AppName string `mapstructure:"APP_NAME"`
	AppPort string `mapstructure:"APP_PORT"`
}

type Wss struct {
	Prefix string `mapstructure:"WSS_PREFIX"`
	Host   string `mapstructure:"WSS_HOST"`
}

type Cookie struct {
	Enable bool   `mapstructure:"COOKIE_ENABLE"`
	Name   string `mapstructure:"COOKIE_NAME"`
	Secret string `mapstructure:"COOKIE_SECRET"`
}

func (w *Wss) Wss() string {
	return w.Host + w.Prefix
}

type Redis struct {
	Enable      bool          `mapstructure:"REDIS_Enable"`
	Host        string        `mapstructure:"REDIS_Host"`
	Passwd      string        `mapstructure:"REDIS_Password"`
	MaxIdle     int           `mapstructure:"REDIS_MaxIdle"`
	MaxActive   int           `mapstructure:"REDIS_MaxActive"`
	IdleTimeout time.Duration `mapstructure:"REDIS_IdleTimeout"`
}

var AppSetting = &App{
	Env:     "development",
	AppName: "go-starter",
	AppPort: "8080",
}

var WssSetting = &Wss{
	Prefix: "",
}

var CookieSetting = &Cookie{
	Enable: false,
	Name:   "_mycookie_name",
	Secret: "_mycookie_secret",
}

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

	// err := godotenv.Load(filename)

	// if err != nil {
	// 	log.Println("Configuration loading " + utils.RootDir() + "/" + filename + " file error, errorMessage=" + fmt.Sprint(err) + ".")
	// 	os.Exit(3)
	// 	return
	// }

	v1 := readConfig(filename, map[string]interface{}{})

	loadAppSettings(v1, filename)
	loadWssSettings(v1, filename)
	loadRedisSettings(v1, filename)
	loadCookieSettings(v1, filename)

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
		log.Println("viper parse RedisConf error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		log.Println("RedisSettings: Enable=" + strconv.FormatBool(RedisConf.Enable) + ", Host=" + RedisConf.Host)
	}

}

func loadWssSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&WssSetting)

	if err != nil {
		log.Println("viper parse WssSetting error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		log.Println("WssSetting: Host=" + WssSetting.Host + ", Prefix=" + WssSetting.Prefix)
	}

}

func loadCookieSettings(v1 *viper.Viper, filename string) {

	err := v1.Unmarshal(&CookieSetting)

	if err != nil {
		log.Println("viper parse CookieSetting error: file=" + filename + " errorMessage=" + fmt.Sprint(err) + ".")
		os.Exit(3)
		return
	} else {
		log.Println("CookieSetting: Enable=" + strconv.FormatBool(CookieSetting.Enable) + ", Name=" + CookieSetting.Name + ", Secret=" + CookieSetting.Secret)
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
	log.Println("viper Configuration loaded " + utils.RootDir() + "/" + filename + " successful.")
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
