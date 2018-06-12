package main

import (
	"configrd"
	"os"
)

func main() {
	var confile = configrd.Config(os.Args[1])
	conf := confile.ReadConfig()

}
