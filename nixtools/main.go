package main

import "log"

func main() {
	if err := runTools(); err != nil {
		log.Fatal(err)
	}
}

func runTools() error {
	if err := update_environment_file(); err != nil {
		return err
	}
	return nil
}
