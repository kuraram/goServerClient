package typefile

import (
	//"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	//"unsafe"
	//"strconv"
)

type FileTrans struct {
	config     Config
	data       []byte // 転送ファイルを格納
	packet_num int
	payloads   map[int][]byte
	IP         string
}

func (ft *FileTrans) OpenYmlFile(filename string) { //YAMLファイルの読み込み

	buf, err := ioutil.ReadFile(filename)
	fmt.Println(string(buf))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, &ft.config)
	if err != nil {
		panic(err)
	}
}

func (ft *FileTrans) OpenTransFile(filename string) { //転送ファイルの読み込み

	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f) // ファイル全体の読み込み
	//fmt.Println(string(ft.data))
	ft.packet_num = int(len(b)) / ft.config.DATA_SIZE
	if int(len(b))%ft.config.DATA_SIZE != 0 {
		ft.packet_num += 1
	}

	ft.data = make([]byte, ft.config.DATA_SIZE*ft.packet_num) // ゼロパディング処理
	for i := 0; i < len(b)-1; i += 1 {
		ft.data[i] = b[i]
	}

	fmt.Printf("File SIZE (byte): %d\n", int(len(b)))
	fmt.Printf("After Zero Padding (byte): %d\n", int(len(ft.data)))
	fmt.Printf("Packet NUM : %d\n", ft.packet_num)

}

func (ft *FileTrans) CreatePayload() { // 先頭4バイトに独自ヘッダを付与

	ft.payloads = map[int][]byte{} //初期化
	for i := 0; i < ft.packet_num; i += 1 {
		data := ft.data[i*ft.config.DATA_SIZE : (i+1)*ft.config.DATA_SIZE]
		header := int_to_bytes(i)
		data = append(header, data...)
		ft.payloads[i] = data
		//fmt.Println(len(ft.payloads[i]))
	}
}

func (ft *FileTrans) SendFile() {

	conn, err := net.Dial("udp", ft.IP)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	for i := 0; i < ft.packet_num; i += 1 {
		fmt.Fprintf(conn, string(ft.payloads[i]))
	}
	conn.Close()

}
