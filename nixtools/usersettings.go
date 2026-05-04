package main

import (
	"errors"
	_ "log"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
)

// Everything will assume the shell is BASH

func getCurrentUser() (string, error) {
	slog.Info("checking for user settings")
	current_user, err := user.Current()
	if err != nil {
		slog.Error("error reading local user", "error", err)
		return "", err
	}
	// Check if running on local machine
	if current_user.Username == "nick" {
		return "", errors.New("current user indicates may be running locally")
	}
	slog.Info("detected username", "username", current_user.Username)
	return current_user.Username, nil
}

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

func (c *Config) updateEnvironmentFile() error {
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
