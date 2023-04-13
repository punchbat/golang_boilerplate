package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	production  string = "production"
	development string = "development"
)

func setEnv(configKey string, envName string) {
	value, exist := os.LookupEnv(envName)
	if exist {
		viper.Set(configKey, value)
	}
}

func setConfigsFromEnv() {
	// set env for db
	setEnv("db.uri", "DB_URI")
	setEnv("db.name", "DB_NAME")
	setEnv("db.options.user", "DB_USER")
	setEnv("db.options.password", "DB_USER_PASSWORD")

	// set env for services
	setEnv("services.email.SMTP_FROM", "SMTP_FROM")
	setEnv("services.email.SMTP_PASSWORD", "SMTP_PASSWORD")

	// set auth env
	setEnv("auth.signing_key", "AUTH_SIGNING_KEY")
	setEnv("auth.token_ttl", "AUTH_TOKEN_TTL")
}

func Init() {
	env, exist := os.LookupEnv("GO_ENV")

	if !exist {
		env = development

		log.Println("GO_ENV is not exist, app will init config at development mode")
	} else {
		log.Printf("App will init config with %s mode", env)
	}

	err := godotenv.Load(fmt.Sprintf(".env.%s", env))
	if err != nil {
		log.Printf("Error loading .env.%s file", env)
	}

	viper.AddConfigPath("./configs") // path to look for the config file in
	viper.SetConfigType("json")      // REQUIRED if the config file does not have the extension in the name

	// cетим конфиги по окружению
	if env == production {
		viper.SetConfigName(production)
	} else {
		viper.SetConfigName(development)
	}

	// cетим конфиги c env файла
	setConfigsFromEnv()
	log.Println("Configs from env setted")

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		log.Fatalf("%s", err.Error())

		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	log.Println("Configs inited")
}
