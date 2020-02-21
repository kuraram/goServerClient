package typefile

import (
	"fmt"
	"io/ioutil"
	"net"
	"unsafe"
)

var DATA_SIZE = 1468

//struct
type FileRetrieve struct {
	data       []byte // 受信ファイルを格納
	datasize   int
	packet_num int
	IP         string
}

func (fr *FileRetrieve) RecvFile() {

	fr.datasize = 22000 //既知

	p := make([]byte, DATA_SIZE)
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	j := 0
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		fr.data = append(fr.data, p...) // 結合
		fmt.Printf("Read a message from %v %s \n", remoteaddr, *(*string)(unsafe.Pointer(&p)))
		p = make([]byte, DATA_SIZE)
		j += 1
		if j == (fr.datasize/DATA_SIZE)+1 { // 全パケットを取得したら終了
			break
		}
	}

	fr.data = fr.data[0 : fr.datasize-1] // ゼロパディング削除
	err = ioutil.WriteFile("message.txt", fr.data, 0755)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

}
