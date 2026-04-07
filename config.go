package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Token  string `toml:"token"`
	ChatID string `toml:"chat_id"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "stt.conf"), nil
}

func loadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("could not read config at %s: %w", configPath, err)
	}

	if cfg.Token == "" || cfg.ChatID == "" {
		return nil, fmt.Errorf("config at %s is missing token or chat_id", configPath)
	}

	return &cfg, nil
}

func setup() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists at %s. Overwrite? (y/N): ", configPath)
		var response string
		_, _ = fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Setup cancelled.")
			return nil
		}
	}

	var cfg Config
	fmt.Print("Enter Telegram Bot Token: ")
	_, _ = fmt.Scanln(&cfg.Token)
	fmt.Print("Enter Telegram Chat ID: ")
	_, _ = fmt.Scanln(&cfg.ChatID)

	if cfg.Token == "" || cfg.ChatID == "" {
		return fmt.Errorf("token and chat_id are required")
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		return err
	}

	fmt.Printf("Config saved to %s\n", configPath)
	return nil
}
