package main

import (
	"fmt"
	"linux-tools/nixtools/pkg"
	"log"
	"log/slog"
)

func main() {
	selection := pkg.DisplayMenu()
	if err := runTools(selection); err != nil {
		log.Fatal(err)
	}
}

func routeRequest(s string) error {
	switch s {
	case "1":
		fmt.Println("option 1")
	case "2":
		fmt.Println("option 2")
	case "3":
		fmt.Println("option 3")
	}
	return nil
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
		slog.Info("option 2 noop")
	case "3":
		slog.Info("option 3 noop")
	default:
		slog.Error("option required")
	}

	return nil
}
