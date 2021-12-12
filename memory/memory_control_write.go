package memory

import (
	"fmt"
)

type writeMemFunc func(uint8) error

// TODO: handle using `interface.cart`. Placeholder for now.
func (m *memory) write_rom(addr uint16) writeMemFunc {
	return func(val uint8) error {
		fmt.Printf("[ROM] Writing to Addr: 0x%04x Val: 0x%02x\n", addr, val)
		return nil
	}
}

func (m *memory) write_vram(addr uint16) writeMemFunc {
	newAddr := addr - VRAM_START
	fmt.Printf("0x%x\n", newAddr)
	return func(val uint8) error {
		m.vRAM[newAddr] = val
		return nil
	}
}

// TODO: implement through cart.
func (m *memory) write_eram(addr uint16) writeMemFunc {
	newAddr := addr - EXT_RAM_START
	return func(val uint8) error {
		m.eRAM[newAddr] = val
		return nil
	}
}

func (m *memory) write_wram(addr uint16) writeMemFunc {
	newAddr := mapwRAMIndex(addr)
	return func(u uint8) error {
		m.wRAM[newAddr] = u
		return nil
	}
}

func (m *memory) write_oam(addr uint16) writeMemFunc {
	newAddr := addr - OAM_START
	return func(u uint8) error {
		m.OAM[newAddr] = u
		return nil
	}
}

func (m *memory) write_io(addr uint16) writeMemFunc {
	newAddr := addr - IO_START
	return func(u uint8) error {
		m.ioRegs[newAddr] = u
		return nil
	}
}

func (m *memory) write_hram(addr uint16) writeMemFunc {
	newAddr := addr - HRAM_START
	return func(u uint8) error {
		m.hRAM[newAddr] = u
		return nil
	}
}

func (m *memory) write_ie(addr uint16) writeMemFunc {
	return func(u uint8) error {
		m.IE[0] = u
		return nil
	}
}
