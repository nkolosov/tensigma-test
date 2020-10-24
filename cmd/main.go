package main

import (
	"github.com/nkolosov/tendigma-test/internal/config"
	"log"
)

func main() {
	log.Printf("start service")

	cfg := config.MustConfigure()
	log.Printf("start service with config %+v", cfg)

	log.Printf("stop service")
}
