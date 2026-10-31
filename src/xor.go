package main

type XorShift32 struct {
	state uint32
}

func NewXorShift32(seed uint32) *XorShift32 {
	if seed == 0 {
		seed = 0x12345678
	}

	return &XorShift32{state: seed}
}

func (this *XorShift32) NextByte() byte {
	this.state ^= this.state << 13
	this.state ^= this.state >> 17
	this.state ^= this.state << 5
	return byte(this.state & 0xFF)
}
