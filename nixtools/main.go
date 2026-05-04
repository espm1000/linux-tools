package main

import "log"

func main() {
	if err := runTools(); err != nil {
		log.Fatal(err)
	}
}

func runTools() error {
	// if _, err := checkOS(); err != nil {
	// 	return err
	// }
	// if err := update_environment_file(); err != nil {
	// 	return err
	// }
	if err := install_apt_dependencies(); err != nil {
		return err
	}
	return nil
}
