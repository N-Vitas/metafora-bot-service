package config

import (
	"encoding/json"
	"os"
)

// Conf Структура конфигурации сервиса
type Conf struct {
	DBUSER string `json:"dbuser"`
	DBPASS string `json:"dbpass"`
	DBHOST string `json:"dbhost"`
	DBPORT int64  `json:"dbport"`
	DBNAME string `json:"dbname"`
}

// New получение конфигурации
func New() (*Conf, error) {
	f, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	conf := &Conf{}
	err = json.NewDecoder(f).Decode(conf)
	return conf, err
}
