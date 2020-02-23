package main

import (
	"./typefile"
)

func main() {

	var fr typefile.FileRetrieve

	filename := "./config.yaml"
	fr.OpenYmlFile(filename)
	fr.Initialize()

	filename = "./retrieval/message.txt"
	fr.RecvPacket()
	fr.JoinPacket(filename)
}
