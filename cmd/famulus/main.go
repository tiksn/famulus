package main

import (
	"log"

	config "github.com/tiksn/famulus/internal/app/famulus"
)

func main() {
	c, err := config.ParseDefaultConfig()
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(c.ListAddress())
}
