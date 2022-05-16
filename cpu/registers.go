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

func GetRegName(idx uint8) string {
	switch idx {
	case A:
		return "A"
	case B:
		return "B"
	case C:
		return "C"
	case D:
		return "D"
	case E:
		return "E"
	case F:
		return "F"
	case H:
		return "H"
	case L:
		return "L"
	}
	return ""
}

const (
	CARRY     uint8 = 0b00010000
	HALFCARRY uint8 = 0b00100000
	NEG       uint8 = 0b01000000
	ZERO      uint8 = 0b10000000
)
