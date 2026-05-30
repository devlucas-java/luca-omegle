package configs

import (
	"fmt"
	"log"

	redis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbName     string `mapstructure:"DB_NAME"`
	DbSSLMode  string `mapstructure:"DB_SSL_MODE"`

	ChHost     string `mapstructure:"CH_HOST"`
	ChPort     string `mapstructure:"CH_PORT"`
	ChUser     string `mapstructure:"CH_USER"`
	ChPassword string `mapstructure:"CH_PASSWORD"`
	ChDB       int    `mapstructure:"CH_DB"`

	ServerPort string `mapstructure:"SERVER_PORT"`
}

var conf *Config

func InitConfig() *Config {
	conf = &Config{}

	viper.SetConfigFile(".default.env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("log in ReadConfig viper, error: %v \n", err.Error())
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		log.Fatalf("log in Unmarshal viper, error: %v \n", err.Error())
	}
	return conf
}

func InitCache(cfg *Config) {

	redis.NewClient(&redis.Options{
		Addr:     cfg.ChHost + ":" + cfg.ChPort,
		Password: cfg.ChPassword,
		DB:       cfg.ChDB,
	})
}

func InitDB() *gorm.DB {

	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.DbHost, conf.DbPort, conf.DbUser, conf.DbPassword, conf.DbName, conf.DbSSLMode)

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("error in init db, %v", err)
	}
	return db
}
