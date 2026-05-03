package main

import "log/slog"

func main() {
	_, err := checkOS()
	if err != nil {
		slog.Info("error checking OS")
	}
}
