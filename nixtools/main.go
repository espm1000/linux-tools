package main

import (
	"fmt"
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
	cu, _ := getCurrentUser()
	hostname, _ := getHostname()
	homeDir, _ := checkEnvironmentFile()

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
	fmt.Println(cfg)
	return nil
}
