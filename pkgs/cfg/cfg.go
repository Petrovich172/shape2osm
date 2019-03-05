package cfg

import (
	"crypto/tls"

	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Configuration - Общая конфигурация сервера.
type Configuration struct {
	PostgresDatabaseCfg PostgresDatabaseCfg `json:"postgres_database_cfg"`
	ServerCfg           struct {
		Port string `json:"port"`
	} `json:"server_cfg"`
	RedisCfg RedisDatabaseCfg `json:"redis_cfg"`
}

// PostgresDatabaseCfg - конфигурация подключений (массив) к разным базам данных Postgres
type PostgresDatabaseCfg []struct {
	Host      string      `json:"host"`
	Port      string      `json:"port"`
	Database  string      `json:"database"`
	User      string      `json:"user"`
	Password  string      `json:"password"`
	EnableTLS *tls.Config `json:"enable_tls"`
}

// RedisDatabaseCfg - конфигурация подключения к БД Redis
type RedisDatabaseCfg struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	Password   string `json:"password"`
	DBIndecies []int  `json:"db_indecies"`
}

// SetParams - инициализация параметров сервера
func (cfg *Configuration) SetParams(fname string) {

	configFile, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		return
	}
	err = json.Unmarshal(configFile, &cfg)
	if err != nil {
		fmt.Println("Error parsing configuration data:", err)
		return
	}

	fmt.Println("Using next configuration:")
	fmt.Println("\t Postgres Database:")
	for i := range (*cfg).PostgresDatabaseCfg {
		fmt.Println("\t\tDB index:", i)
		fmt.Println("\t\tHost:", (*cfg).PostgresDatabaseCfg[i].Host)
		fmt.Println("\t\tPort:", (*cfg).PostgresDatabaseCfg[i].Port)
		fmt.Println("\t\tDatabase name:", (*cfg).PostgresDatabaseCfg[i].Database)
		fmt.Println("\t\tUser", (*cfg).PostgresDatabaseCfg[i].User)
		fmt.Println("\t\tPassword", (*cfg).PostgresDatabaseCfg[i].Password)
		fmt.Println("\t\tEnable TLS:", (*cfg).PostgresDatabaseCfg[i].EnableTLS)
	}
	fmt.Println("\tServer:")
	fmt.Println("\t\tServer port:", (*cfg).ServerCfg.Port)
}
