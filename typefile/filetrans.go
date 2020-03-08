package typefile

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	//"unsafe"
	"strconv"
)

type FileTrans struct {
	config               Config
	info                 Info
	b                    []byte
	data                 []byte           // 転送ファイルを格納
	packet_num           int              // 全パケット数
	packet_num_per_block int              // 1ブロックのパケット数
	payloads             map[int][]byte   // 全ペイロードを格納
	sockets              map[int]net.Conn // Port番号に対応したソケット
	IP                   string
}

func (ft *FileTrans) OpenYmlFile(filename string) { //YAMLファイルの読み込み

	buf, err := ioutil.ReadFile(filename)
	//fmt.Println(string(buf))
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
		panic(err)
	}
	defer f.Close()

	ft.b, err = ioutil.ReadAll(f) // ファイル全体の読み込み
	//fmt.Println(string(ft.data))

	/*
		ft.data = make([]byte, ft.config.DATA_SIZE*ft.packet_num) // ゼロパディング処理
		for i := 0; i < len(b)-1; i += 1 {
			ft.data[i] = b[i]
		}
	*/

	fmt.Printf("File SIZE (byte): %d\n", int(len(ft.b)))
	//fmt.Printf("After Zero Padding (byte): %d\n", int(len(ft.data)))
	//fmt.Printf("Packet NUM : %d\n", ft.packet_num)

}

func (ft *FileTrans) ReadInfo(payload string) { // OFCからのペイロードの読み込み

	err := json.Unmarshal([]byte(payload), &ft.info)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Phase NUM : %d\n", ft.info.PhaseNum)
	fmt.Printf("Split NUM : %d\n", ft.info.SplitNum)

	ft.packet_num = int(len(ft.b)) / ft.config.DATA_SIZE // ブロック数1の場合のパケット数
	if int(len(ft.b))%ft.config.DATA_SIZE != 0 {
		ft.packet_num += 1
	}

	// 各ブロックのパケット数
	ft.packet_num_per_block = ft.packet_num / ft.info.SplitNum
	if ft.packet_num%ft.info.SplitNum != 0 {
		ft.packet_num_per_block += 1
	}

	ft.packet_num = ft.info.SplitNum * ft.packet_num_per_block
	ft.data = make([]byte, ft.config.DATA_SIZE*ft.packet_num) // ゼロパディング処理

	for i := 0; i < len(ft.b)-1; i += 1 { // 送信データを代入
		ft.data[i] = ft.b[i]
	}

	fmt.Printf("Packet NUM : %d\n", ft.packet_num)
	fmt.Printf("Packet NUM per Block : %d\n", ft.packet_num_per_block)

}

func (ft *FileTrans) CreatePayload() { // 先頭4バイトに独自ヘッダを付与

	var tool Tool
	ft.payloads = map[int][]byte{} //初期化

	for i := 0; i < ft.packet_num; i += 1 {
		data := ft.data[i*ft.config.DATA_SIZE : (i+1)*ft.config.DATA_SIZE]
		header := tool.int_to_bytes(i)
		data = append(header, data...)
		ft.payloads[i] = data
		//fmt.Println(header)
	}
}

func (ft *FileTrans) CreateSockets() { // 転送に使用するソケットを作成

	ft.sockets = map[int]net.Conn{}

	for phase := 0; phase < ft.info.PhaseNum; phase += 1 {
		for _, block := range ft.info.Blocks[phase] {
			port := 10000 + 100*phase + block
			addr := ft.info.MulticastIP + ":" + strconv.Itoa(port)
			conn, err := net.Dial("udp", addr)
			if err != nil {
				panic(err)
			}
			fmt.Println(port)
			ft.sockets[port] = conn
		}
	}

}

func (ft *FileTrans) SendPacket(port int, pos int) {
	fmt.Fprintf(ft.sockets[port], string(ft.payloads[pos]))
}

func (ft *FileTrans) Packet_num_per_block() int {
	return ft.packet_num_per_block
}

func (ft *FileTrans) Phase_num() int {
	return ft.info.PhaseNum
}

func (ft *FileTrans) Blocks() [][]int {
	return ft.info.Blocks
}
