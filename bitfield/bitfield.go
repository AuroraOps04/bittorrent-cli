package bitfield

type Bitfield []byte

func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8 //in which byte the bit is
	offset := index % 8
	// 判断 第 byteIndex 个 byte 的 offset 位是否为 1
	// offset 从 0 开始计数
	return bf[byteIndex]>>(7-offset)&1 != 0
}

func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8 //in which byte the bit is
	offset := index % 8
	// 将 第 byteIndex 个 byte 的 offset 位置为 1
	bf[byteIndex] |= 1 << (7 - offset)
}
