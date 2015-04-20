package main

import (
	"github.com/tumdum/gosmbclient"
	"log"
	"os"
)

func rmdir() {
	url := os.Args[2]
	if err := gosmbclient.RmDir(url); err != nil {
		log.Fatalf("Failed to remove dir '%s': %s", url, err)
	}
}
