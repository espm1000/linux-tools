package main

import (
	"log"
)

type Config struct {
	currentUser    string
	hostname       string
	homeDiretory   string
	packageManager string
	distro         string
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

	return Config{
		currentUser:  cu,
		hostname:     hostname,
		homeDiretory: homeDir,
	}, nil

}

func runTools() error {
	cfg, err := generateConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := cfg.updateEnvironmentFile(); err != nil {
		return err
	}
	return nil
}
