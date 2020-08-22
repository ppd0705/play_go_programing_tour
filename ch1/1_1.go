package main

import (
	"flag"
	"log"
)

func main1() {
	var name string
	flag.StringVar(&name, "name", "default value", "Help msg")
	flag.StringVar(&name, "n", "default value", "Help msg")
	flag.Parse()
	log.Printf("name: %s\n", name)
}
