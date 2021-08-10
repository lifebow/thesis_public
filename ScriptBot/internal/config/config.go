package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

type WorkerConfig struct {
	RabbitmqHost        string
	RabbitmqDefaultUser string
	RabbitmqDefaultPass string
	ScriptCheckFolder   string
	UrlLink             string
	Jwt 				string
}

func MustReadWorkerConfig() WorkerConfig {
	var cfg WorkerConfig
	err:= godotenv.Load()
	if err!=nil{
		log.Warn().Msg("Error while reading .env config")
	}
	cfg.RabbitmqHost=os.Getenv("RabbitmqHost")
	cfg.RabbitmqDefaultUser=os.Getenv("RabbitmqUser")
	cfg.RabbitmqDefaultPass=os.Getenv("RabbitmqPass")
	cfg.ScriptCheckFolder=os.Getenv("ScriptCheckFolder")
	cfg.UrlLink=os.Getenv("UrlLink")
	cfg.Jwt=os.Getenv("Jwt")
	return cfg
}
