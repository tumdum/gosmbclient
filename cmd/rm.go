package main

import (
	"github.com/tumdum/gosmbclient"
	"log"
	"os"
)

func rm() {
	url := os.Args[2]
	if err := gosmbclient.Unlink(url); err != nil {
		log.Fatalf("Failed to unlink '%s': %s", url, err)
	}
}
