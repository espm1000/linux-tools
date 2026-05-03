package main

import (
	"bufio"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

func checkOS() (string, error) {
	os := runtime.GOOS
	switch os {
	case "linux":
		slog.Info("operating system detected", "os", os)
		distro, err := getLinuxDistro()
		if err != nil {
			return "", err
		}
		slog.Info("distro detected", "distro", distro)
	default:
		slog.Info("failed to detect OS")
	}

	return os, nil
}

func getLinuxDistro() (string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		s, found := strings.CutPrefix(line, "PRETTY_NAME=")
		if found {
			return s, nil
		}
	}
	return "", nil
}
