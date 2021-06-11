package diablosutils

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppName      string `mapstructure:"APP_NAME"`
	AllowOrigins string `mapstructure:"ALLOW_ORIGINS"`
	DbHost       string `mapstructure:"DB_HOST"`
	DbPort       string `mapstructure:"DB_PORT"`
	DbUser       string `mapstructure:"DB_USER"`
	DbPassword   string `mapstructure:"DB_PASSWORD"`
	DbName       string `mapstructure:"DB_NAME"`
	RbUser       string `mapstructure:"RB_USER"`
	RbPassword   string `mapstructure:"RB_PASSWORD"`
	RbHost       string `mapstructure:"RB_HOST"`
	RbPort       string `mapstructure:"RB_PORT"`
	RdPassword   string `mapstructure:"REDIS_PASSWORD"`
	RdHost       string `mapstructure:"REDIS_HOST"`
	RdPort       string `mapstructure:"REDIS_PORT"`
	AccessKey    string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshKey   string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AppPort      string `mapstructure:"APP_PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
