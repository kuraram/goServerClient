package typefile

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"sync"
	//"os"
	//"strconv"
	//"unsafe"
)

//struct
type FileRetrieve struct {
	config               Config
	info                 Info
	data                 []byte // 受信ファイルを格納
	packet_num           int
	packet_num_per_block int                  // 1ブロックのパケット数
	payloads             map[int][]byte       // 全ペイロードを格納
	sockets              map[int]*net.UDPConn // Port番号に対応したソケット
	IP                   string
	count                int
	ret                  chan int
	res                  map[int]chan int
	sig                  map[int]chan int
	mux                  sync.Mutex
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

func (fr *FileRetrieve) ReadInfo(payload string) { // OFCからのペイロードの読み込み

	err := json.Unmarshal([]byte(payload), &fr.info)
	if err != nil {
		panic(err)
	}

	fmt.Printf("File SIZE (byte): %d\n", fr.info.DataSize)
	fmt.Printf("Phase NUM : %d\n", fr.info.PhaseNum)
	fmt.Printf("Split NUM : %d\n", fr.info.SplitNum)

}

func (fr *FileRetrieve) Initialize() { // 諸々の情報の計算

	fr.packet_num = fr.info.DataSize / fr.config.DATA_SIZE
	if fr.info.DataSize%fr.config.DATA_SIZE != 0 {
		fr.packet_num += 1
	}
	fmt.Printf("Packet NUM : %d\n", fr.packet_num)

	// 各ブロックのパケット数
	fr.packet_num_per_block = fr.packet_num / fr.info.SplitNum
	if fr.packet_num%fr.info.SplitNum != 0 {
		fr.packet_num_per_block += 1
	}

	fmt.Printf("Packet NUM per Block : %d\n", fr.packet_num_per_block)

	fr.payloads = map[int][]byte{} //ペイロード部分の初期化

}

func (fr *FileRetrieve) CreateSockets() { // 受信に使用するソケットを作成

	fr.sockets = map[int]*net.UDPConn{}

	for phase := 0; phase < fr.info.PhaseNum; phase += 1 {
		for _, block := range fr.info.Blocks[phase] {
			port := 10000 + 100*phase + block
			addr := net.UDPAddr{
				Port: port,
				IP:   net.ParseIP(fr.info.MulticastIP),
			}
			conn, err := net.ListenMulticastUDP("udp", nil, &addr)
			if err != nil {
				panic(err)
			}
			fr.sockets[port] = conn
		}
	}

}

func (fr *FileRetrieve) SetChanMap() {

	fr.res = map[int]chan int{}
	fr.sig = map[int]chan int{}

	ports := fr.GetKeysPortNUM()
	for _, port := range ports {
		fr.res[port] = make(chan int, fr.packet_num_per_block)
		fr.sig[port] = make(chan int, fr.packet_num_per_block)
	}

	fr.ret = make(chan int, len(ports))
}

func (fr *FileRetrieve) RetrievePacket(port int) { // パケットの取得

	//fmt.Println(port)

	for i := 0; i < fr.packet_num_per_block; i++ {

		p := make([]byte, fr.config.PAYLOAD_SIZE)
		conn := fr.sockets[port]
		_, _, err := conn.ReadFromUDP(p)
		if err != nil {
			panic(err)
		}
		go fr.GetData(p, port)

	}
	for {
		//fmt.Println(len(sig))
		if len(fr.sig[port]) >= fr.packet_num_per_block {
			break
		}
	}

	fr.ret <- port
}

func (fr *FileRetrieve) GetData(p []byte, port int) { // 独自ヘッダとペイロード部分の取得

	fr.res[port] <- 1

	var tool Tool
	num := tool.bytes_to_int(p[:fr.config.CUSTOM_HEAD_SIZE])

	// 排他制御
	fr.mux.Lock()
	fr.payloads[num] = p[fr.config.CUSTOM_HEAD_SIZE:]
	fr.mux.Unlock()

	//fmt.Println(num)
	fr.count += 1
	fr.sig[port] <- num

	<-fr.res[port]

}

func (fr *FileRetrieve) JoinPacket(filename string) { // 全パケットの結合と書き込み

	//fmt.Println(fr.payloads[1])
	//for i := 0; i < fr.packet_num; i += 1 {
	for i := 0; i < fr.packet_num; i += 1 {
		//fmt.Println(i)
		fr.data = append(fr.data, fr.payloads[i]...)
	}

	fr.data = fr.data[:fr.info.DataSize-1] // ゼロパディング削除

	err := ioutil.WriteFile(filename, fr.data, 0755)
	if err != nil {
		panic(err)
	}

}

func (fr *FileRetrieve) GetKeysPortNUM() []int { // Socketと結び付けられたポート番号を取得
	keys := make([]int, len(fr.sockets))
	i := 0
	for k := range fr.sockets {
		keys[i] = k
		i++
	}
	return keys
}

func (fr *FileRetrieve) Count() int {
	return fr.count
}

func (fr *FileRetrieve) Packet_Num() int {
	return fr.packet_num
}

func (fr *FileRetrieve) Ret() chan int {
	return fr.ret
}
