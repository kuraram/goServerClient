package main

import (
	"./typefile"
	//"fmt"
)

func main() {

	var ft typefile.FileTrans

	filename := "./config.yaml"
	ft.OpenYmlFile(filename)

	filename = "./transfer/message"
	ft.OpenTransFile(filename)

	payload_json := `{
		"multicast_ip":"239.0.0.1",
		"data_size":22000,
		"split_num":4,
		"coded_bum":0,
		"phase_num":1,
		"blocks":[[0,1,2,3]]
	}`
	ft.ReadInfo(payload_json)

	ft.CreatePayload()
	ft.CreateSockets()

	// パケット送信部分
	Blocks := ft.Blocks()
	for phase := 0; phase < ft.Phase_num(); phase += 1 { // フェーズ毎

		for num := 0; num < ft.Packet_num_per_block(); num += 1 {

			for _, block := range Blocks[phase] { // ブロック毎

				pos := block*ft.Packet_num_per_block() + num
				port := 10000 + 100*phase + block
				ft.SendPacket(port, pos)
				//fmt.Println(port, pos)

			}
		}
	}

}
