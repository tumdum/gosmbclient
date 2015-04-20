package main

import (
	"fmt"
	"github.com/tumdum/gosmbclient"
	"os"
)

func die(e error) {
	fmt.Fprintf(os.Stderr, "%s\n", e.Error())
	os.Exit(1)
}

func usage() {
	fmt.Printf("%s command [options]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	gosmbclient.Init(SimpleAuth, 0)
	if len(os.Args) <= 1 {
		usage()
	}
	switch os.Args[1] {
	case "ls":
		ls()
	case "mkdir":
		mkdir()
	case "rmdir":
		rmdir()
	case "cat":
		cat()
	case "cp":
		cp()
	case "rm":
		rm()
	case "mv":
		mv()
	}
}
