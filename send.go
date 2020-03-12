package main

import (
	"./typefile"
	"fmt"
	"time"
)

func main() {

	var ft typefile.FileTrans

	filename := "./config.yaml"
	ft.OpenYmlFile(filename)

	filename = "./transfer/dummy"
	ft.OpenTransFile(filename)

	payload_json := `{
		"multicast_ip":"239.0.0.1",
		"data_size":100000000,
		"split_num":6,
		"coded_bum":0,
		"phase_num":1,
		"blocks":[[0,1,2,3,4,5]]
	}`
	ft.ReadInfo(payload_json)

	ft.CreatePayload()
	ft.CreateSockets()
	ft.SetChan()

	// パケット送信部分
	PhaseBlocks := ft.Blocks()
	packet_num := ft.Packet_num()
	phase_num := ft.Phase_num()
	packet_num_per_block := ft.Packet_num_per_block()
	//fmt.Println(Blocks)

	start := time.Now()
	phase := 0
	for { // フェーズ毎

		Blocks := PhaseBlocks[phase]
		pos := 0
		fmt.Printf("\nPhase %d :　\n", phase+1)
		for _, block := range Blocks { // ブロック毎
			fmt.Printf(" %d ", block)
		}

		for {

			if pos >= packet_num_per_block {
				break
			}

			for _, block := range Blocks { // ブロック毎

				port := 10000 + 100*phase + block
				//fmt.Println(port, block, pos)
				go ft.SendPacket(port, block, pos)

			}

			pos += 1
		}
		for {
			//fmt.Println(len(ft.Comp()))
			if len(ft.Comp()) >= packet_num {
				//fmt.Println(len(fr.Ret()))
				break
			}
		}

		phase += 1
		if phase == phase_num {
			break
		}

	}

	end := time.Now()
	fmt.Printf("\n%f秒\n", (end.Sub(start)).Seconds())

	ft.CloseSockets()

}
