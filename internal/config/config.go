package config

import "os"

type Config struct {
	Init bool
	Addr string
	Env  string
}

type DataBaseConfig struct {
    User        string
    Password    string
    Port        string
    DBName      string
}


func Load() *Config {
	addr := os.Getenv("ADDR")
	
	if addr == "" {
		addr = ":8080"
	}

	env := os.Getenv("ENV")
	
	if env == "" {
		env = "dev"
	}


	return &Config{Addr: addr, Env: env, Init: true}
}
