package typefile

type Config struct {
	DATA_SIZE        int `yaml:"data_size"`
	PAYLOAD_SIZE     int `yaml:"payload_size"`
	CUSTOM_HEAD_SIZE int `yaml:"custom_head_size"`
}
