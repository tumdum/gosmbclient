package main

import (
	"github.com/tumdum/gosmbclient"
	"io/ioutil"
	"log"
	"os"
)

func cat() {
	url := os.Args[2]
	f, err := gosmbclient.Open(url, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Opening file '%v' failed: %s", url, err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Reading file '%v' failed: %s", url, err)
	}
	os.Stdout.Write(b)
}
