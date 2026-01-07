package main

import (
	"fmt"
	"log"

	"github.com/thomas-reed/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	if err = cfg.SetUser("treed"); err != nil {
		log.Fatalf("Error setting user in config: %v", err)
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	fmt.Println(cfg.DbURL)
	fmt.Println(cfg.CurrentUsername)
}