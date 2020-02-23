package typefile

type Config struct {
	DATA_SIZE      int `yaml:"data_size"`
	PAYLOAD_SIZE   int `yaml:"payload_size"`
	UNIQ_HEAD_SIZE int `yaml:"uniq_head_size"`
}
