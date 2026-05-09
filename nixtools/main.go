package main

import (
	"linux-tools/nixtools/pkg"
	"log"
	"log/slog"
)

func main() {
	selection, err := pkg.DisplayMenu()
	if err != nil {
		log.Fatal(err)
	}
	if err := runTools(selection); err != nil {
		slog.Error("error occured", "error", err)
		log.Fatal(err)
	}
}

func runTools(s string) error {
	cfg, err := pkg.GenerateConfig()
	if err != nil {
		return err
	}
	switch s {
	case "1":
		if err := pkg.UpdateEnvironmentFile(cfg); err != nil {
			return err
		}
		if err := pkg.InstallAptDependencies(cfg); err != nil {
			return err
		}
	case "2":
		if err = pkg.InstallDevTools(cfg, false); err != nil {
			return err
		}
	case "3":
		if err := pkg.InstallDockerDependencies(cfg); err != nil {
			return err
		}
	default:
		pkg.Exit()
	}

	return nil
}
