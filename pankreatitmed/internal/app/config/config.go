package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int

	JWT   JWTConfig   `mapstructure:"jwt"`
	Redis RedisConfig `mapstructure:"redis"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Issuer string        `mapstructure:"issuer"`
	TTL    time.Duration `mapstructure:"ttl"`
}

type RedisConfig struct {
	Addr     string        `mapstructure:"addr"`
	DB       int           `mapstructure:"db"`
	Password string        `mapstructure:"password"`
	TTL      time.Duration `mapstructure:"ttl"`
}

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	log.Infof("jwt.ttl raw (string): %q", viper.GetString("jwt.ttl"))
	log.Infof("jwt map: %#v", viper.GetStringMap("jwt"))
	log.Infof("all settings: %#v", viper.AllSettings())
	if err != nil {
		return nil, err
	}

	cfg := &Config{}           // создаем объект конфига
	err = viper.Unmarshal(cfg) // читаем информацию из файла,
	// конвертируем и затем кладем в нашу переменную cfg
	if err != nil {
		return nil, err
	}
	//cfg := &Config{}
	//if err := viper.Unmarshal(cfg,
	//	viper.DecodeHook(mapstructure.StringToTimeDurationHookFunc()),
	//); err != nil {
	//	return nil, err
	//}

	//cfg.JWT.TTL = viper.GetDuration("jwt_ttl")
	//cfg.Redis.TTL = viper.GetDuration("redis_ttl")
	log.Info("config parsed")

	return cfg, nil
}
