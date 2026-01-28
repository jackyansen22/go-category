package config

import "github.com/spf13/viper"

type Config struct {
	AppPort string
	DBUrl   string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	return &Config{
		AppPort: viper.GetString("APP_PORT"),
		DBUrl: "postgres://" +
			viper.GetString("DB_USER") + ":" +
			viper.GetString("DB_PASSWORD") + "@" +
			viper.GetString("DB_HOST") + ":" +
			viper.GetString("DB_PORT") + "/" +
			viper.GetString("DB_NAME") + "?sslmode=require",
	}
}
