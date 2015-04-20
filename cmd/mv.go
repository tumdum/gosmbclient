package main

import (
	"github.com/tumdum/gosmbclient"
	"log"
	"os"
)

func mv() {
	src := os.Args[2]
	dst := os.Args[3]
	if err := gosmbclient.Rename(src, dst); err != nil {
		log.Fatalf("Failed to move from '%s' to '%s': %s", src, dst, err)
	}
}
