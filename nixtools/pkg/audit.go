package pkg

import (
	"bufio"
	"errors"
	"log/slog"
	"os"
	"os/user"
	"runtime"
	"strings"
)

type Config struct {
	currentUser    string
	hostname       string
	homeDiretory   string
	packageManager string
	distro         string
	os             string
}

func GenerateConfig() (*Config, error) {
	cu, err := getCurrentUser()
	if err != nil {
		slog.Error("error getting user", "error", err)
		return nil, err
	}

	hostname, err := getHostname()
	if err != nil {
		return nil, err
	}

	homeDir, err := checkEnvironmentFile(cu)
	if err != nil {
		return nil, err
	}

	cfg, err := checkOS()
	if err != nil {
		return nil, err
	}

	return &Config{
		currentUser:    cu,
		hostname:       hostname,
		homeDiretory:   homeDir,
		os:             cfg.os,
		distro:         cfg.distro,
		packageManager: cfg.packageManager,
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

func checkOS() (*Config, error) {
	cfg := &Config{}
	os := runtime.GOOS
	switch os {
	case "linux":
		slog.Info("operating system detected", "os", os)
		config, err := getLinuxDistro()
		if err != nil {
			return nil, err
		}
		cfg = config
	default:
		slog.Info("failed to detect OS")
	}
	return &Config{
		os:             os,
		distro:         cfg.distro,
		packageManager: cfg.packageManager,
	}, nil
}

func getLinuxDistro() (*Config, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, err
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
				return &Config{
					packageManager: "apt",
					distro:         "debian",
				}, nil
			case strings.Contains(normalizeString(line), "redhat") || strings.Contains(normalizeString(line), "fedora"):
				return &Config{
					packageManager: "dnf",
					distro:         "redhat",
				}, nil
			default:
				slog.Error("unable to detect distribution", "response", line)
				return nil, errors.New("no distro detected")
			}
		}
	}
	return nil, nil
}

func normalizeString(s string) string {
	return strings.ToLower(s)
}
