package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB     Postgres `mapstructure:"postgres"`
	Server Server   `mapstructure:"server"`
}

type Server struct {
	Port string `mapstructure:"port"`
}

type Postgres struct {
	Host     string `mapstructure:"db_host"`
	DBPort   string `mapstructure:"db_port"`
	Username string `mapstructure:"db_username"`
	DBName   string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"db_sslmode"`
	Password string `mapstructure:"db_password"`
}

func LoadConfig() Config {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("config file not found")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("unable to decode into structure, %v", err)
	}

	log.Printf("Loaded config: %+v", config)
	return config
}
