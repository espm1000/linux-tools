package main

import (
	"bufio"
	"errors"
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
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("error closing io stream", "error", err)
		}
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		_, found := strings.CutPrefix(line, "PRETTY_NAME=")

		if found {
			switch {
			case strings.Contains(normalizeString(line), "debian") || strings.Contains(normalizeString(line), "ubuntu"):
				return "debian", nil
			case strings.Contains(normalizeString(line), "redhat") || strings.Contains(normalizeString(line), "fedora"):
				return "redhat", nil
			default:
				slog.Error("unable to detect distribution", "response", line)
				return "", errors.New("no distro detected")
			}
		}
	}
	return "", nil
}

func normalizeString(s string) string {
	return strings.ToLower(s)
}
