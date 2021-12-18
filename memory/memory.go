package memory

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

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
	W_RAM_E_MID = 0xDDFF
	W_RAM_END   = 0xDFFF

	ECHO_RAM_START = 0xE000
	ECHO_RAM_END   = 0xFDFF

	OAM_START = 0xFE00
	OAM_END   = 0xFE9F

	HRAM_START = 0xFF80
	HRAM_END   = 0xFFFE

	IO_START = 0xFF00
	IO_END   = 0xFF7F

	BOOT_LOADER_FLAG = 0xFF50

	INTERRUPT_ENABLE = 0xFFFF
)

func mapwRAMIndex(addr uint16) uint16 {
	var newAddr uint16
	if addr >= W_RAM_START && addr <= W_RAM_END {
		newAddr = addr - W_RAM_START
	} else if addr >= ECHO_RAM_START && addr <= ECHO_RAM_END {
		// Map Echo ram to the same range and work RAM range
		newAddr = (addr-ECHO_RAM_START)/(ECHO_RAM_END-ECHO_RAM_START)*(W_RAM_E_MID-W_RAM_START) + W_RAM_START
		newAddr -= W_RAM_START
	}
	return newAddr
}

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
	OAM     []uint8 // 0xFE00 - 0xFE9F
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
		hRAM:       make([]uint8, 126),
		romData:    make([]uint8, 32*1024),
		ioRegs:     make([]uint8, 128),
		OAM:        make([]uint8, 159),
		bootloader: bootLoader,
		rom:        ROM,
		IE:         make([]uint8, 1),
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

func (m *memory) getReadMemBlock(addr uint16) readMemFunc {
	if m.isBootLoaderLoaded() && addr < 0x100 {
		return m.read_boot_loader(addr)
	} else if addr <= ROM_END {
		return m.read_rom_data(addr)
	} else if addr <= VRAM_END {
		return m.read_vram_data(addr)
	} else if addr <= EXT_RAM_END {
		return m.read_ext_ram(addr)
	} else if addr <= W_RAM_END {
		return m.read_wram(addr)
	} else if addr <= ECHO_RAM_END {
		return m.read_wram(addr)
	} else if addr <= OAM_END {
		return m.read_oam(addr)
	} else if addr <= IO_END {
		return m.read_io(addr)
	} else if addr <= HRAM_END {
		return m.read_hram(addr)
	} else if addr == INTERRUPT_ENABLE {
		return m.read_IE()
	}
	return nil
}

func (m *memory) getWriteMemBlock(addr uint16) writeMemFunc {
	if addr <= ROM_END {
		return m.write_rom(addr)
	} else if addr <= VRAM_END {
		return m.write_vram(addr)
	} else if addr <= EXT_RAM_END {
		return m.write_eram(addr)
	} else if addr <= W_RAM_END {
		return m.write_wram(addr)
	} else if addr <= ECHO_RAM_END {
		return m.write_wram(addr)
	} else if addr <= OAM_END {
		return m.write_oam(addr)
	} else if addr <= IO_END {
		return m.write_io(addr)
	} else if addr <= HRAM_END {
		return m.write_hram(addr)
	} else if addr == INTERRUPT_ENABLE {
		return m.write_ie(addr)
	}
	return nil
}

func (m *memory) MemRead(addr uint16) *uint8 {
	// fmt.Printf("Reading from 0x%04x\n", addr)
	mem := m.getReadMemBlock(addr)
	if mem == nil {
		// handle error
		log.Fatalf("Received invalid memory address range: 0x%04X\n", addr)
		return nil
	}
	return mem()
}

func (m *memory) MemWrite(addr uint16, val uint8) error {
	// fmt.Printf("Writing to: 0x%04x\n", addr)
	mem := m.getWriteMemBlock(addr)
	return mem(val)
}
