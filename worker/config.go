package worker

import (
	"encoding/json"
	"io/ioutil"
)

type Server struct {
	Ip     string `json:"ip"`
	Port   string `json:"port"`
	ApiKey string `json:"api_key"`
}

type Rest struct {
	DBUrl  string `json:"db_url"`
	ApiKey string `json:"api_key"`
}

type DBPath struct {
	Main          string `json:"main"`
	ServerState   string `json:"server_state"`
	MapsAvailable string `json:"maps_available"`
}

type Firebase struct {
	Rest   `json:"rest"`
	DBPath `json:"db_path"`
}

type Mongo struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type configData struct {
	Server   `json:"server"`
	Firebase `json:"firebase"`
	Mongo    `json:"mongo"`
}

// Configure Use it with a config.json file to configure the server properties
func Configure() configData {
	newConfig := new(configData)
	data, err := ioutil.ReadFile("config.json")
	Check(err)
	json.Unmarshal(data, newConfig)
	return *newConfig
}
