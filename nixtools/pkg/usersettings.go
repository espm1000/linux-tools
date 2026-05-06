package pkg

import (
	_ "log"
	"log/slog"
	"os"
	"path/filepath"
)

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
	slog.Info("updating user settings file", "path", c.homeDiretory)
	template, err := os.ReadFile("./internal/templates/bashrc.template")
	if err != nil {
		return err
	}
	slog.Info("reading template file")
	f, err := os.OpenFile(filepath.Join(c.homeDiretory, ".bashrc"), os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("error closing file", "error", err)
		}
	}()
	sTemplate := string(template)
	slog.Debug("template contents", "contents", sTemplate)
	if _, err := f.WriteString(sTemplate); err != nil {
		slog.Error("error writing to environment file", "error", err)
		return err
	}
	return nil
}

func getHostname() (string, error) {
	return os.Hostname()
}
