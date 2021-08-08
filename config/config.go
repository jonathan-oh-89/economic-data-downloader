package config

type Configuration struct {
	DBConfig *DBConfig
}

type DBConfig struct {
	HostName string  `json:"hostName"`
	DbName   string  `json:"dbName"`
	UserName string  `json:"userName"`
	Password string  `json:"password"`
	Port     float64 `json:"port"`
}

var Config = &Configuration{}
