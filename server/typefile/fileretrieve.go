package typefile

import (
	//"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"unsafe"
)

var DATA_SIZE = 1468

//struct
type FileRetrieve struct {
	data       []byte // 転送ファイルを格納
	packet_num int
	IP         string
}

func (fr *FileRetrieve) sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func (fr *FileRetrieve) RecvFile() {

	p := make([]byte, 1468) //2048Bytes
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		fmt.Printf("Read a message from %v %s \n", remoteaddr, *(*string)(unsafe.Pointer(&p)))
		err = ioutil.WriteFile("hello-world.txt", p, 0755)
		if err != nil {
			fmt.Printf("Some error %v\n", err)
			return
		}

		go fr.sendResponse(ser, remoteaddr)
	}

}
