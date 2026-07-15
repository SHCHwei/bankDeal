package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsDotEnvExample(t *testing.T) {
	
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env.example"), []byte("addr=8080\nenv=dev\nuser=user\npassword=user1\nport=3306\ndatabase=bankDeal\ndatabase_host=localhost\n"), 0o644); err != nil {
		t.Fatalf("write .env.example: %v", err)
	}

	data := loadEnvExampleFile(filepath.Join(dir, ".env.example"))

	for key, val := range data {
		if val == "" {
			t.Errorf("config setting [%s] is empty", key)
		}
	}



	addr := getEnvValue(data, "addr", ":8080")
	env := getEnvValue(data, "env", "dev")
	user := getEnvValue(data, "user", "user")
	password := getEnvValue(data, "password", "user1")
	port := getEnvValue(data, "port", "3306")
	dbName := getEnvValue(data, "database", "bankDeal")
	database_host := getEnvValue(data, "database_host", "localhost")


	if addr == "" {
		t.Errorf("config setting [address] is error")
	}
	if env == "" {
		t.Errorf("config setting [env] is error")
	}
	if user == "" {
		t.Errorf("config setting [user] is error")
	}
	if password == "" {
		t.Errorf("config setting [password] is error")
	}
	if port == "" {
		t.Errorf("config setting [port] is error")
	}
	if dbName == "" {
		t.Errorf("config setting [database] is error")
	}
	if database_host == "" {
		t.Errorf("config setting [database_host] is error")
	}

}
