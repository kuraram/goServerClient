package typefile

type Tool struct {
	bs []byte
	is int
}

func (t *Tool) int_to_bytes(in int) []byte { // int -> 4 Bytes

	t.bs = []byte{}
	t.bs = []byte{
		byte(in >> 24 & 0xFF),
		byte(in >> 16 & 0xFF),
		byte(in >> 8 & 0xFF),
		byte(in & 0xFF),
	}

	return t.bs
}

func (t *Tool) bytes_to_int(in []byte) int { // 4Bytes -> int

	t.is = int((int32(in[0])<<24 | int32(in[1])<<16 | int32(in[2])<<8 | int32(in[3])))

	return t.is
}
