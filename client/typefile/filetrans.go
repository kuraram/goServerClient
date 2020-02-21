package typefile

import (
	//"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

var DATA_SIZE = 1468

//struct
type FileTrans struct {
	data       []byte // 転送ファイルを格納
	packet_num int
	IP         string
}

func (ft *FileTrans) OpenTransFile(filename string) { //転送ファイルの読み込み

	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error")
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f) // ファイル全体の読み込み
	//fmt.Println(string(ft.data))
	ft.packet_num = int(len(b)) / DATA_SIZE
	if int(len(b))%DATA_SIZE != 0 {
		ft.packet_num += 1
	}

	ft.data = make([]byte, DATA_SIZE*ft.packet_num) // ゼロパディング処理
	for i := 0; i < len(b)-1; i++ {
		ft.data[i] = b[i]
	}

	fmt.Printf("File SIZE (byte): %d\n", int(len(b)))
	fmt.Printf("After Zero Padding (byte): %d\n", int(len(ft.data)))
	fmt.Printf("Packet NUM : %d\n", ft.packet_num)

}

func (ft *FileTrans) SendFile() {

	//p := make([]byte, 2048)
	p := ft.data[0:1467]
	fmt.Println(string(p))
	conn, err := net.Dial("udp", ft.IP)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, string(p))
	/*
		p = make([]byte, 2048)
		_, err = bufio.NewReader(conn).Read(p)
		if err == nil {
			fmt.Printf("%s\n", p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
	*/
	conn.Close()

}
