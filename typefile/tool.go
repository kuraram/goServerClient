package typefile

type Tool struct {
	bs []byte
}

func (t *Tool) int_to_bytes(x int) []byte { // int -> 4 Bytes

	t.bs = []byte{}
	t.bs = []byte{
		byte(x >> 24 & 0xFF),
		byte(x >> 16 & 0xFF),
		byte(x >> 8 & 0xFF),
		byte(x & 0xFF),
	}

	return t.bs
}
