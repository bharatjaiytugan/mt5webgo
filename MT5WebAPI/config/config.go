package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// MT5Config holds the connection parameters loaded from .env file.
type MT5Config struct {
	Server   string
	Port     int
	Login    uint64
	Password string
}

// Load reads the .env file and returns MT5Config.
// It searches for .env in the current directory, then parent directories up to 3 levels.
func Load() (*MT5Config, error) {
	envPath := findEnvFile()
	if envPath == "" {
		return nil, fmt.Errorf(".env file not found")
	}

	vars, err := parseEnvFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse .env: %v", err)
	}

	cfg := &MT5Config{}

	if v, ok := vars["MT5_SERVER"]; ok {
		cfg.Server = v
	} else {
		return nil, fmt.Errorf("MT5_SERVER not set in .env")
	}

	if v, ok := vars["MT5_PORT"]; ok {
		port, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid MT5_PORT: %v", err)
		}
		cfg.Port = port
	} else {
		cfg.Port = 443
	}

	if v, ok := vars["MT5_LOGIN"]; ok {
		login, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid MT5_LOGIN: %v", err)
		}
		cfg.Login = login
	} else {
		return nil, fmt.Errorf("MT5_LOGIN not set in .env")
	}

	if v, ok := vars["MT5_PASSWORD"]; ok {
		cfg.Password = v
	} else {
		return nil, fmt.Errorf("MT5_PASSWORD not set in .env")
	}

	return cfg, nil
}

func findEnvFile() string {
	dir, _ := os.Getwd()
	for i := 0; i < 4; i++ {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			return path
		}
		dir = filepath.Dir(dir)
	}
	return ""
}

func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	vars := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		vars[key] = val
	}
	return vars, scanner.Err()
}
