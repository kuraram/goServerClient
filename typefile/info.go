package typefile

type Info struct {
	MulticastIP string  `json:"multicast_ip"`
	DataSize    int     `json:"data_size"`
	SplitNum    int     `json:"split_num"`
	CodedNum    int     `json:"coded_num"`
	PhaseNum    int     `json:"phase_num"`
	Blocks      [][]int `json:"blocks"`
}
