package config

import (
	"log"

	"github.com/spf13/viper"
)

// type Config struct {
// 	DB   Postgres
// 	Port string `mapstructure:"PORT"`
// }

type Config struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	DBName   string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
	Password string `mapstructure:"DB_PASSWORD"`
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
