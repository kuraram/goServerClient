package main

import (
	"./typefile"
	"fmt"
)

func main() {

	var fr typefile.FileRetrieve

	filename := "./config.yaml"
	fr.OpenYmlFile(filename)

	payload_json := `{
		"multicast_ip":"239.0.0.1",
		"data_size":22000,
		"split_num":4,
		"coded_bum":0,
		"phase_num":1,
		"blocks":[[0,1,2,3]]
	}`
	fr.ReadInfo(payload_json)
	fr.Initialize()

	fr.CreateSockets()
	fr.SetChanMap()

	ports := fr.GetKeysPortNUM()
	//packet_num := fr.Packet_Num()
	/*
		countdown := 0
		fmt.Println(packet_num)
		for countdown < packet_num {
			fmt.Println(countdown)
			for _, port := range ports {
				fmt.Println(port)
				go fr.RetrievePacket(port)
				countdown = fr.Count()
			}
		}
	*/
	fmt.Println(ports)
	for _, port := range ports {
		go fr.RetrievePacket(port)
	}
	for {
		if len(fr.Comp()) >= len(ports) {
			//fmt.Println(len(fr.Ret()))
			break
		}
	}

	//fmt.Println(fr.Count())
	filename = "./retrieval/message.txt"
	fr.JoinPacket(filename)
}
