package main

import (
	"flag"
	"fmt"
	"github.com/tumdum/gosmbclient"
	"log"
	"os"
)

func ls() {
	long := flag.Bool("l", false, "use a long listing format")
	flag.CommandLine.Parse(os.Args[2:])
	fmt.Println("long", *long)
	url := os.Args[len(os.Args)-1]
	d, err := gosmbclient.OpenDir(url)
	if err != nil {
		log.Fatalf("Failed to list contents of '%v': %s", url, err)
	}
	defer d.Close()
	names, err := d.List()
	if err != nil {
		die(err)
	}
	for _, n := range names {
		if *long {
			info, err := gosmbclient.Stat(n)
			if err == nil {
				printFileInfo(info)
			}
		} else {
			fmt.Println(n)
		}
	}
	os.Exit(0)
}

func printFileInfo(info os.FileInfo) {
	fmt.Printf("%-60s %10d %s %v %v\n", info.Name(), info.Size(), info.Mode(),
		info.ModTime(), info.IsDir())
}
