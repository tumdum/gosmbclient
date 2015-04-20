package main

import (
	"github.com/tumdum/gosmbclient"
	"io"
	"log"
	"os"
)

func cp() {
	// TODO(klak): use same permisions in dst as source
	src := os.Args[2]
	dst := os.Args[3]
	fdst, err := gosmbclient.Create(dst, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to open file '%v': %s", dst, err)
	}
	defer fdst.Close()
	fsrc, err := os.Open(src)
	if err != nil {
		log.Fatalf("Failed to open file '%v': %s", src, err)
	}
	fsrc.Close()
	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		log.Fatalf("Failed to copy: %s", err)
	}
}
