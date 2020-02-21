package main

import (
	"./typefile"
)

func main() {

	var ft typefile.FileTrans
	filename := "./message"
	ft.OpenTransFile(filename)

	ft.IP = "127.0.0.1:1234"
	ft.SendFile()
}
