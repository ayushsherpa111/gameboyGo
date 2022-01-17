package cartridge

import "fmt"

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
	title         string
	cgbFlag       byte
	cartridgeType string
	sgbFlag       bool
	romSize       byte
	ramSize       int
	isJap         bool
	versionNum    byte
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
	c.sgbFlag = romData[sgb_flag] == 0x03
	c.romSize = romData[rom_size]
	c.ramSize = ramSizes[romData[ram_size]]
	c.isJap = romData[lang_code] == 0x00
	c.versionNum = romData[version_num]
}

func NewCart(romData []byte) *cart {
	c := &cart{}
	c.parseHeader(romData)
	return c
}
