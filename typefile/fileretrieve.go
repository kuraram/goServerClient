package typefile

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	//"unsafe"
)

//struct
type FileRetrieve struct {
	config     Config
	data       []byte // 受信ファイルを格納
	datasize   int
	packet_num int
	payloads   map[int][]byte // 全ペイロードを格納
	IP         string
}

func (fr *FileRetrieve) OpenYmlFile(filename string) { //YAMLファイルの読み込み

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, &fr.config)
	if err != nil {
		panic(err)
	}
}

func (fr *FileRetrieve) Initialize() {

	fr.datasize = 22000 //既知

	fr.packet_num = fr.datasize / fr.config.DATA_SIZE
	if fr.datasize%fr.config.DATA_SIZE != 0 {
		fr.packet_num += 1
	}
	fmt.Printf("Packet NUM : %d\n", fr.packet_num)

}

func (fr *FileRetrieve) RecvPacket() {

	var tool Tool
	fr.payloads = map[int][]byte{} //初期化

	p := make([]byte, fr.config.PAYLOAD_SIZE)
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
	for { // パケット受信部分
		//_, remoteaddr, err := ser.ReadFromUDP(p)
		_, _, err := ser.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}

		fr.payloads[tool.bytes_to_int(p[:4])] = p[4:]
		//fmt.Printf("Read a message from %v %s \n", remoteaddr, *(*string)(unsafe.Pointer(&p)))
		p = make([]byte, fr.config.PAYLOAD_SIZE)
		j += 1
		if j == (fr.datasize/fr.config.DATA_SIZE)+1 { // 全パケットを取得したら終了
			break
		}
	}
}

func (fr *FileRetrieve) JoinPacket(filename string) { // 全パケットの結合

	//fmt.Println(fr.payloads[1])
	for i := 0; i < fr.packet_num; i += 1 {
		fr.data = append(fr.data, fr.payloads[i]...)
	}

	fr.data = fr.data[:fr.datasize-1] // ゼロパディング削除

	err := ioutil.WriteFile(filename, fr.data, 0755)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

}
