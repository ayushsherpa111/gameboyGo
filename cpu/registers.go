package cpu

type register uint8

const (
	A = 0
	B
	C
	D
	E
	F
	H
	L
)

const (
	CARRY     = 0b00010000
	HALFCARRY = 0b00100000
	SUB       = 0b01000000
	ZERO      = 0b10000000
)

func initRegisters() [8]register {
	return [8]register{}
}
