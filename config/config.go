package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var filename string = ".env"

// LoadEnvVars will load a ".env[.development|.test]" file if it exists and set ENV vars.
// Useful in development and test modes. Not used in production.
func init() {

	var env string = os.Getenv("GIN_ENV")

	filename = ".env." + env
	err := godotenv.Load(filename)

	if err != nil {

		log.Println("loading " + filename + " file error, use .env file, errorMessage=" + fmt.Sprint(err) + ".")
		filename = ".env"

		err = godotenv.Load(filename)

		if err != nil {
			log.Println("load .env file error: errorMessage=" + fmt.Sprint(err) + ".")
			os.Exit(3)
			return
		}

	}

	log.Println(filename + " loaded.")

}

// GetEnv returns an environment variable or a default value if not present
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
