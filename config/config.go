package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	MainPath   string `json:"main_path"`
	ServerPort uint16 `json:"server_port"`
}

func (c Config) String() string {
	return fmt.Sprintf("Загруженный канфиг:\n  --> Порт сервера: %v\n  --> Каталог: %v\n", c.ServerPort, c.MainPath)
}

func MustLoad() *Config {
	path := "config/config.json"
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("cannot open config file: %v", err))
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		panic(fmt.Sprintf("cannot decode config file: %v", err))
	}

	fmt.Println(cfg)

	return &cfg
}
