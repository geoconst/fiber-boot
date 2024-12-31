package app

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port       int    `json:"port"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
}

func NewConfig() *Config {
	var config Config
	data, _ := os.ReadFile("./config.json")
	err := json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
