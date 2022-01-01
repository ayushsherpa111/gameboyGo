package cpu

const (
	A = iota
	B
	C
	D
	E
	F
	H
	L
)

const (
	CARRY     uint8 = 0b00010000
	HALFCARRY uint8 = 0b00100000
	NEG       uint8 = 0b01000000
	ZERO      uint8 = 0b10000000
)
