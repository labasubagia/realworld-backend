package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	PostgresSource       string `mapstructure:"POSTGRES_SOURCE"`
	PostgresMigrationURL string `mapstructure:"POSTGRES_MIGRATION_URL"`
	HTTPServerAddress    string `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey    string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
