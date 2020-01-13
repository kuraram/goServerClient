package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Server is Running at localhost:8888")
	conn, _ := net.ListenPacket("udp", "localhost:8888")
	defer conn.Close()

	buffer := make([]byte, 1500)
	for {
		// 通信読込 + 接続相手アドレス情報が受取
		length, remoteAddr, _ := conn.ReadFrom(buffer)
		fmt.Printf("Received from %v: %v\n", remoteAddr, string(buffer[:length]))
		conn.WriteTo([]byte("Hello, World !"), remoteAddr)
	}

}
