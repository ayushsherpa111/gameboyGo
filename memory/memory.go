package memory

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ayushsherpa111/gameboyEMU/interfaces"
)

type Mem interface {
	MemRead(addr uint16) *uint8
	MemWrite(addr uint16, val uint8) error
}

const (
	ROM_END = 0x7FFF

	VRAM_START = 0x8000
	VRAM_END   = 0x9FFF

	EXT_RAM_START = 0xA000
	EXT_RAM_END   = 0xBFFF

	W_RAM_START = 0xC000
	W_RAM_END   = 0xDFFF

	IO_START = 0xFF00
	IO_END   = 0xFF7F

	BOOT_LOADER_FLAG = 0xFF50

	INTERRUPT_ENABLE = 0xFFFF
)

// 0x0000 - 0x3FFF : ROM Bank 0 [Must be handled thru memory bank controller]
// 0x4000 - 0x7FFF : ROM Bank 1 - Switchable [Must be handled thru memory bank controller]
// 0x8000 - 0x9FFF : 8KB VRAM
// 0xA000 - 0xBFFF : Cartridge RAM [Must be handled thru memory bank controller]
// 0xC000 - 0xCFFF : 4KB Work RAM
// 0xD000 - 0xDFFF : 4KB Work RAM
// 0xE000 - 0xFDFF : Reserved - Echo RAM [Mirror of 0xC000 - 0xDDFF]
// 0xFE00 - 0xFE9F : Object Attribute Memory
// 0xFEA0 - 0xFEFF : Reserved - Unusable
// 0xFF00 - 0xFF7F : I/O Registers
// 0xFF80 - 0xFFFE : High RAM
// 0xFFFF - 0xFFFF : Interrupt Enable Flag

type memory struct {
	bootloader []uint8
	rom        string

	romData []uint8 // 0x0000 - 0x8000
	vRAM    []uint8 // 0x8000 - 0x9FFF
	eRAM    []uint8 // 0xA000 - 0xBFFF
	wRAM    []uint8 // 0xC000 - 0xDFFF
	hRAM    []uint8 // 0xFF80 - 0xFFFE
	ioRegs  []uint8 // 0xFF00 - 0xFF70
	IE      []uint8 // 0xFFFF

	// romswp []uint8
	cart interfaces.Cart
}

func (m *memory) ReadIO(addr uint16) uint8 {
	addr -= IO_START
	return m.ioRegs[addr]
}

func (m *memory) isBootLoaderLoaded() bool {
	if m.ReadIO(BOOT_LOADER_FLAG) == 0 {
		return true
	}
	return false
}

func InitMem(bootLoader []byte, ROM string) (*memory, error) {
	fmt.Println(len(bootLoader))
	if len(bootLoader) != 256 {
		return nil, errors.New("Invalid Bootloader provided")
	}

	mem := &memory{
		vRAM:       make([]uint8, 8*1024),
		eRAM:       make([]uint8, 8*1024),
		wRAM:       make([]uint8, 8*1024),
		romData:    make([]uint8, 32*1024),
		ioRegs:     make([]uint8, 128),
		bootloader: bootLoader,
		rom:        ROM,
		IE:         make([]uint8, 1),
		// romswp:     make([]uint8, 256),
	}

	if e := mem.loadROM(); e != nil {
		return nil, e
	}

	return mem, nil
}

func (m *memory) loadROM() error {
	fmt.Println("Loading ROM file...")
	romData, err := ioutil.ReadFile(m.rom)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to load ROM file. Reason: %s", err.Error()))
	}
	if len(romData) > (32 << 10) {
		return errors.New("Invalid ROM length")
	}
	copy(m.romData, romData)
	fmt.Println("ROM file Succesfully loaded.")
	return nil
}

func (m *memory) getMemBlock(addr uint16) ([]uint8, uint16) {
	var newAddrMap uint16
	var memBlock []uint8

	if m.isBootLoaderLoaded() && addr < 0x100 {
		newAddrMap = addr
		memBlock = m.bootloader
	} else if addr <= ROM_END {
		newAddrMap = addr
		memBlock = m.romData
	} else if addr <= VRAM_END {
		newAddrMap = addr - VRAM_START
		memBlock = m.vRAM
	} else if addr <= EXT_RAM_END {
		newAddrMap = addr - EXT_RAM_START
		memBlock = m.eRAM
	} else if addr <= W_RAM_END {
		newAddrMap = addr - W_RAM_START
		memBlock = m.wRAM
	} else if addr <= IO_END {
		newAddrMap = addr - IO_START
		memBlock = m.ioRegs
	} else if addr == INTERRUPT_ENABLE {
		newAddrMap = 0
		memBlock = m.IE
	}

	return memBlock, newAddrMap
}

func (m *memory) MemRead(addr uint16) *uint8 {
	mem, newAddr := m.getMemBlock(addr)
	return &mem[newAddr]
}

func (m *memory) MemWrite(addr uint16, val uint8) error {
	fmt.Printf("Writing to: 0x%04x\n", addr)
	mem, newAddr := m.getMemBlock(addr)
	mem[newAddr] = val
	return nil
}

// func (m *memory) GetByte(addr uint16) (error, *byte) {
// 	return nil, &m.memory[addr]
// }

// func (m *memory) SetByte(addr uint16, val byte) error {
// 	m.memory[addr] = val
// 	return nil
// }
