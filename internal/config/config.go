package config

import (
	"bufio"
	"os"
	"strings"
	"log"
)

type Config struct {
	Addr string
	Env  string
}

type DataBaseConfig struct {
	User     string
	Password string
	Port     string
	DBName   string
	Host     string
}



func Load() (*Config, *DataBaseConfig) {
	values := loadEnvExampleFile(".env.example")

	addr := getEnvValue(values, "addr", ":8080")
	env := getEnvValue(values, "env", "dev")

	user := getEnvValue(values, "user", "user")
	password := getEnvValue(values, "password", "user1")
	port := getEnvValue(values, "port", "3306")
	dbName := getEnvValue(values, "database", "bankDeal")
	database_host := getEnvValue(values, "database_host", "localhost")


	return &Config{
		Addr: addr,
		Env:  env,
	}, &DataBaseConfig{
		User:     user,
		Password: password,
		Port:     port,
		DBName:   dbName,
		Host:     database_host,
	}

}

func loadEnvExampleFile(path string) map[string]string {
	values := make(map[string]string)

	file, err := os.Open(path)

	if err != nil {
		log.Fatal("Failed to open env file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
 
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		values[key] = value
	}

	if scanner.Err() != nil {
		log.Fatal("Failed to read env file:", scanner.Err())
	}


	return values
}



func getEnvValue(values map[string]string, targetKey string, defaultValue string) string {

	if value, exists := values[targetKey]; exists {
		return value
	}
	return defaultValue
}