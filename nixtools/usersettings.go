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

func getConfig() (*Config, error) {
	current_user, _ := getCurrentUser()
	user_home_path, err := checkEnvironmentFile()
	if err != nil {
		slog.Error("error getting configuration details")
		return nil, err
	}
	hostname, _ := getHostname()
	return &Config{
		currentUser:  current_user,
		homeDiretory: user_home_path,
		hostname:     hostname,
	}, nil

}

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

func checkEnvironmentFile() (string, error) {
	cu, err := getCurrentUser()
	if err != nil {
		return "", err
	}
	slog.Info("detected current user", "user", cu)
	var userPath = "/home/" + cu
	slog.Info("assuming user path", "path", userPath)
	if _, err := os.Stat(filepath.Join(userPath, ".bashrc")); err == nil {
		slog.Info("bash environment file found")
		return userPath, nil
	} else {
		slog.Info("bashrc file not found, creating")
		return "", err
	}
}

func updateEnvironmentFile() error {
	slog.Info("checking for environment file(s)")
	path, err := checkEnvironmentFile()
	if err != nil {
		return err
	}
	slog.Info("updating user settings file", "path", path)
	template, err := os.ReadFile("./internal/templates/bashrc.template")
	if err != nil {
		return err
	}
	slog.Info("reading template file")
	f, err := os.OpenFile(filepath.Join(path, ".bashrc"), os.O_WRONLY|os.O_APPEND, 0644)
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
