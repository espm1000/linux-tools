package main

import (
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
)

// Everything will assume the shell is BASH

func getCurrentUser() (string, error) {
	current_user, err := user.Current()
	if err != nil {
		return "", err
	}
	slog.Info("detected username", "username", current_user.Username)
	return current_user.Username, nil
}
func check_for_environment_file() (string, error) {
	cu, err := getCurrentUser()
	if err != nil {
		return "", err
	}
	var userPath = "/home/" + cu
	if _, err := os.Stat(userPath + "/.bashrc"); err == nil {
		slog.Info("bash environment file found")
		return userPath, nil
	} else {
		slog.Info("bashrc file not found, creating")
		return "", err
	}
}

func update_environment_file() error {
	path, err := check_for_environment_file()
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
