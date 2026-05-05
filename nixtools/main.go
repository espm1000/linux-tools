package main

import (
	"log"
	"os"
)

type Config struct {
	currentUser    string
	hostname       string
	homeDiretory   string
	packageManager string
	distro         string
	os             string
}

func main() {
	if err := runTools(); err != nil {
		log.Fatal(err)
	}
}

func generateConfig() (Config, error) {
	cu, err := getCurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	hostname, _ := getHostname()
	homeDir, _ := checkEnvironmentFile(cu)
	cfg, _ := checkOS()

	return Config{
		currentUser:    cu,
		hostname:       hostname,
		homeDiretory:   homeDir,
		os:             cfg.os,
		distro:         cfg.distro,
		packageManager: cfg.packageManager,
	}, nil

}

func runTools() error {
	verbose := os.Getenv("NIX_VERBOSE")
	if verbose == "" {
		verbose = "false"
	}
	cfg, err := generateConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := cfg.updateEnvironmentFile(); err != nil {
		return err
	}
	if err := cfg.installAptDependencies(verbose); err != nil {
		return err
	}
	return nil
}
