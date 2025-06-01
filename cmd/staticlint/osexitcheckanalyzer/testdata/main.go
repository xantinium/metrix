package main

import "os"

func main() {
	os.Exit(1) // want "main func calls to os.Exit"
}
