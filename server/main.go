package main

import (
	"./typefile"
)

func main() {

	var fr typefile.FileRetrieve

	filename := "../config.yaml"
	fr.OpenYmlFile(filename)

	fr.RecvFile()
}
