package main

import (
	"fmt"
	"log"

	config "github.com/tiksn/famulus/internal/app/famulus"
)

func main() {
	c, err := config.ParseDefaultConfig()
	if err != nil {
		log.Fatalln(err)
	}

	adrs, err2 := c.ListAddress()
	if err2 != nil {
		log.Fatalln(err2)
	}
	for _, n := range adrs {
		fmt.Println(c.GetAddress(n))
	}

	path, err3 := c.GetContactsCsvFilePath()
	if err3 != nil {
		log.Fatalln(err3)
	}
	fmt.Println(path)
}
