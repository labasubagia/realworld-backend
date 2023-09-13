package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	PostgresSource       string `mapstructure:"POSTGRES_SOURCE"`
	PostgresMigrationURL string `mapstructure:"POSTGRES_MIGRATION_URL"`

	MongoSource string `mapstructure:"MONGO_SOURCE"`

	HTTPServerPort int `mapstructure:"HTTP_SERVER_PORT"`

	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`

	TestRepo string `mapstructure:"TEST_REPO"`
}

func (c Config) IsTestAllRepo() bool {
	return c.TestRepo == "all"
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
