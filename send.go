package main

import (
	"./typefile"
)

func main() {

	var ft typefile.FileTrans

	filename := "./config.yaml"
	ft.OpenYmlFile(filename)

	filename = "./transfer/message"
	ft.OpenTransFile(filename)

	ft.CreatePayload()

	ft.IP = "127.0.0.1:1234"
	ft.SendFile()
}
