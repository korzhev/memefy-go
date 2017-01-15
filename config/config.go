package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Debug		bool
	Token		string
	TmpPath		string
}

func GetConf() (*Configuration) {
	file, errF := os.Open("config.json")
	if errF != nil {
		log.Panic(errF)
	}
	decoder := json.NewDecoder(file)
	conf := new(Configuration)
	err := decoder.Decode(&conf)
	if err != nil {
		log.Panic(err)
	}
	return conf
}

