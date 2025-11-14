package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)
const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config,error){
	var cfg Config
	home , err := os.UserHomeDir()
	if err != nil { return cfg,err }

	path := filepath.Join(home, configFileName)
    dat, err := os.ReadFile(path)
	if err != nil { return cfg,err }

	if err := json.Unmarshal(dat, &cfg); err != nil {
        return cfg,err
    }

	return cfg,nil
}

//func write(cfg Config) error{
//
//}

func (cfg *Config) SetUser(name string)  error{
	cfg.CurrentUserName = name
	cfg_updated, err := json.Marshal(cfg)
    if err != nil {
        return err
    }
	home , err := os.UserHomeDir()
	if err != nil { return err }

	path := filepath.Join(home, configFileName)
	err = os.WriteFile(path, cfg_updated,0644)
	if err != nil { return err }
	
	return nil
}