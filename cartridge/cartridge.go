package cartridge

import (
	"fmt"
	"log"

	mbc "github.com/ayushsherpa111/gameboyEMU/cartridge/MBC"
	"github.com/ayushsherpa111/gameboyEMU/interfaces"
)

const (
	title_start = 0x134
	title_end   = 0x143

	cgb_flag = 0x143

	license_code_start = 0x144
	license_code_end   = 0x145

	sgb_flag       = 0x146
	cartridge_type = 0x147

	rom_size = 0x148
	ram_size = 0x149

	lang_code = 0x14A

	version_num = 0x14C
)

type cart struct {
	title          string
	cgbFlag        byte
	cartridgeType  string
	cartridgeTypeN uint8
	sgbFlag        bool
	romSize        byte
	ramSize        int
	isJap          bool
	versionNum     byte
	rom0           []uint8
	romn           []uint8
	currentRomBank uint8
	mode           bool
	mbc            interfaces.MBC
}

func (c *cart) HeaderInfo() {
	fmt.Printf("Title: %s\n", c.title)
	fmt.Printf("CGB Flag : %02xh\n", c.cgbFlag)
	fmt.Printf("SGB support : %v\n", c.sgbFlag)
	fmt.Printf("Cartridge Type: %s\n", c.cartridgeType)
	fmt.Printf("ROM Size: %vKB\n", 32<<c.romSize)
	fmt.Printf("RAM Size: %v\n", c.ramSize)
	fmt.Printf("Is Japanese: %v\n", c.isJap)
	fmt.Printf("Version Number: %v\n", c.versionNum)
}

// TODO: Implement Read/Write for different MBC
func (c *cart) parseHeader(romData []byte) {
	c.title = string(romData[title_start:title_end])
	c.cgbFlag = romData[cgb_flag]
	c.cartridgeType = cartridgeTypes[romData[cartridge_type]]
	c.cartridgeTypeN = romData[cartridge_type]
	c.sgbFlag = romData[sgb_flag] == 0x03
	c.romSize = romData[rom_size]
	c.ramSize = ramSizes[romData[ram_size]]
	c.isJap = romData[lang_code] == 0x00
	c.versionNum = romData[version_num]
}

func (c *cart) parseRom(romData []byte) {
	// Non swapable Rom area
	copy(c.rom0, romData[:16<<10])

	numberOfBanks := ((32 << c.romSize) / 16) - 2
	// Contains the remaining banks for swap with ROM1
	if c.mbc = mbc.NewMBC(romData[16<<10:], c.cartridgeTypeN, uint8(numberOfBanks)); c.mbc == nil {
		log.Fatalln("Invalid MBC code receivedd")
	}
}

// addr is going to be in the range of 0x0000 to 0x7FFF
func (c *cart) ReadROM(addr uint16) *uint8 {
	if addr <= mbc.ROM0_END {
		return &c.rom0[addr]
	}
	// TODO: return value from romn after getting the slice value from current MBC
	return &c.romn[addr-mbc.ROMN_START]
}

func (c *cart) WriteROM(addr uint16, val uint8) {
	// change memory when addr between 0x2000 and 0x4000
	if addr >= 0x2000 && addr < 0x4000 {
		c.currentRomBank = (val & 0x1F) + 1
		c.romn = c.mbc.GetSlice(c.currentRomBank)
	}
}

func NewCart(romData []byte) *cart {
	c := &cart{
		rom0:           make([]uint8, 16<<10),
		currentRomBank: 1,
	}
	c.parseHeader(romData)
	c.parseRom(romData)
	c.romn = c.mbc.GetSlice(c.currentRomBank)
	return c
}
