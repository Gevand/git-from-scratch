package utils

func CombineBytes(b1, b2 byte) byte {
	return b1<<4 | b2
}

func PackHexaDecimal(bytes []byte) []byte {
	return_bytes := make([]byte, len(bytes)/2)
	for i := 0; i < len(return_bytes); i++ {
		new_byte := CombineBytes(bytes[i*2], bytes[i*2+1])
		return_bytes[i] = new_byte
	}
	return return_bytes
}
