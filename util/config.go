package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBSource          string `mapstructure:"DB_SOURCE"`
	DBMigrationURL    string `mapstructure:"DB_MIGRATION_URL"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
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
