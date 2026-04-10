package utils

import (
	"os"
	"path/filepath"
)

const tokenDir = ".opentunnel"
const tokenFile = "token"

func GetTokenPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, tokenDir, tokenFile), nil
}

func SaveToken(token string) error {
	path, err := GetTokenPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(token), 0600)
}

func LoadToken() (string, error) {
	// 1. Check Env
	if token := os.Getenv("OPENTUNNEL_TOKEN"); token != "" {
		return token, nil
	}

	// 2. Check File
	path, err := GetTokenPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func TokenExists() bool {
	path, err := GetTokenPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}
