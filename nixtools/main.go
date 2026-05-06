package main

import (
	"linux-tools/nixtools/pkg"
	"log"
)

func main() {
	if err := runTools(); err != nil {
		log.Fatal(err)
	}
}

func runTools() error {
	// verbose := os.Getenv("NIX_VERBOSE")
	cfg, err := pkg.GenerateConfig()
	if err != nil {
		return err
	}
	if err := pkg.UpdateEnvironmentFile(cfg); err != nil {
		return err
	}
	return nil
}
