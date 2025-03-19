package config

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		DB     Postgres `mapstructure:"postgres"`
		Server Server   `mapstructure:"server"`
	}

	Server struct {
		Port string `mapstructure:"port"`
	}

	Postgres struct {
		Host           string `mapstructure:"db_host"`
		DBPort         string `mapstructure:"db_port"`
		DBExternalPort string `mapstructure:"db_external_port"`
		Username       string `mapstructure:"db_username"`
		DBName         string `mapstructure:"db_name"`
		SSLMode        string `mapstructure:"db_sslmode"`
		Password       string `mapstructure:"db_password"`
	}
)

func LoadConfig() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	configData, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			key := "${" + pair[0] + "}"
			value := pair[1]
			configData = bytes.ReplaceAll(configData, []byte(key), []byte(value))
		}
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadConfig(bytes.NewReader(configData)); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into structure, %v", err)
	}

	log.Printf("Loaded config: %+v", config)
	return config
}
