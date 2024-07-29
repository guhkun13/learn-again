package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const AppModeDev = "dev"
const AppModePreview = "pre"
const AppModeProduction = "prod"

func LoadEnv() (env *EnvironmentVariable, err error) {
	log.Info().Msg("Load Env Here")

	envFile := ".env"
	v := viper.New()

	// read static env file
	v.SetConfigFile(envFile)
	err = v.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Str("filename", envFile).Msg("viper error read config")
	}

	v.AutomaticEnv()
	err = v.Unmarshal(&env)
	if err != nil {
		log.Error().Err(err).Msg("viper error unmarshal config")
	}

	return
}
