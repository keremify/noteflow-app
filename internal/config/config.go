package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app"`
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"database"`
	JWT      JWT      `mapstructure:"jwt"`
}

type App struct {
	Name string `mapstructure:"name"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Database struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWT struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int    `mapstructure:"expires_in"` // in hours
}

var Configurations Config

func Get() Config {

	return Configurations
}

func Set() error {

	_ = godotenv.Load(".env")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("internal/config")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found or error reading it: %v", err)
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&Configurations); err != nil {
		return fmt.Errorf("error occured during unmarshal")
	}
	return nil
}
