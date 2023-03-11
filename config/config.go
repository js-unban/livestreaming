package config

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	DB_NAME              string
	TOKEN_DURATION_HOURS int32
	ISSUER               string
}

func GetConfig(params ...string) Configuration {
	configuration := Configuration{}
	env := "dev"
	if len(params) > 0 {
		env = params[0]
	}
	fileName := fmt.Sprintf("./config/%s.json", env)
	gonfig.GetConf(fileName, &configuration)
	return configuration
}
