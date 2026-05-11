package pkg

import (
	_ "log"
	"log/slog"
	"os"
	"path/filepath"
)

type UserConfig struct {
	Config Config
}

// Everything will assume the shell is BASH

func checkEnvironmentFile(user string) (string, error) {
	var userPath = "/home/" + user
	slog.Info("assuming user path", "path", userPath)
	if _, err := os.Stat(filepath.Join(userPath, ".bashrc")); err == nil {
		slog.Info("bash environment file found")
		return userPath, nil
	} else {
		slog.Info("bashrc file not found, creating")
		if _, err := os.Create(filepath.Join(userPath, ".bashrc")); err != nil {
			return "", err
		}
		slog.Info("done")
		return "", err
	}
}

func UpdateEnvironmentFile(c *Config) error {
	if err := GenerateTemplates(c); err != nil {
		return err
	}
	return nil
}

func getHostname() (string, error) {
	return os.Hostname()
}
