package config

import (
	"log"
	"os"

	"github.com/magiconair/properties"
)

type Config struct {
	JWTSecret string `properties:"jwt.secret"`
	JWTApiKey string `properties:"jwt.apikey"`
}

func GetConfig() Config {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	p := properties.MustLoadFile(path+"/application.properties", properties.UTF8)
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
