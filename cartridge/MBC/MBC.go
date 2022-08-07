package mbc

import "github.com/ayushsherpa111/gameboyEMU/interfaces"

// memory address consts
const (
	ROM0_START = 0x0000
	ROM0_END   = 0x3FFF

	ROMN_START = 0x4000
	ROMN_END   = 0x7FFF
)

func NewMBC(romdata []uint8, cartType uint8, bankCount uint8) interfaces.MBC {
	switch cartType {
	case 0:
		return &romOnly{
			bank: romdata,
		}
	case 0x1:
		fallthrough
	case 0x2:
		fallthrough
	case 0x3:
		return &mbc1{
			bank:      romdata,
			bankCount: bankCount,
		}
	case 0x5:
		fallthrough
	case 0x6:
		return &mbc2{}
	case 0x0F:
		fallthrough
	case 0x10:
		fallthrough
	case 0x11:
		fallthrough
	case 0x12:
		fallthrough
	case 0x13:
		return &mbc3{
			bank:      romdata,
			bankCount: bankCount,
		}
	default:
		return nil
	}
}
