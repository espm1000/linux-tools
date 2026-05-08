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
	switch s {
	case "1":
		cfg, err := pkg.GenerateConfig()
		if err != nil {
			return err
		}
		if err := pkg.UpdateEnvironmentFile(cfg); err != nil {
			return err
		}
		if err := pkg.InstallAptDependencies(cfg); err != nil {
			return err
		}
	case "2":
		cfg, err := pkg.GenerateConfig()
		if err != nil {
			return err
		}
		if err = pkg.InstallDevTools(cfg, false); err != nil {
			return err
		}
	case "3":
		slog.Info("option 3 noop")
	default:
		pkg.Exit()
	}

	return nil
}
