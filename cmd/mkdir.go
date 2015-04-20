package main

import (
	"github.com/tumdum/gosmbclient"
	"log"
	"os"
)

func mkdir() {
	url := os.Args[2]
	if err := gosmbclient.MkDir(url, os.ModeDir); err != nil {
		log.Fatalf("Failed to create dir '%s': %s", url, err)
	}
}
