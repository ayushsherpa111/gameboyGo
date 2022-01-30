package memory

type readMemFunc func() *uint8

func (m *memory) read_io(addr uint16) readMemFunc {
	newAddr := addr - IO_START
	return func() *uint8 {
		return &m.ioRegs[newAddr]
	}
}

func (m *memory) read_boot_loader(addr uint16) readMemFunc {
	return func() *uint8 {
		return &m.bootloader[addr]
	}
}

func (m *memory) read_rom_data(addr uint16) readMemFunc {
	return func() *uint8 {
		return &m.romData[addr]
	}
}

func (m *memory) read_vram_data(addr uint16) readMemFunc {
	newAddr := addr - VRAM_START
	return func() *uint8 {
		return &m.vRAM[newAddr]
	}
}

// TODO: EXT RAM from Cartridge
func (m *memory) read_ext_ram(addr uint16) readMemFunc {
	newAddr := addr - EXT_RAM_START
	return func() *uint8 {
		return &m.eRAM[newAddr]
	}
}

func (m *memory) read_wram(addr uint16) readMemFunc {
	newAddr := mapwRAMIndex(addr)
	// fmt.Printf("New addr for WRAM: 0x%x\n", newAddr)
	return func() *uint8 {
		return &m.wRAM[newAddr]
	}
}

func (m *memory) read_oam(addr uint16) readMemFunc {
	newAddr := addr - OAM_START
	return func() *uint8 {
		return &m.OAM[newAddr]
	}
}

func (m *memory) read_hram(addr uint16) readMemFunc {
	newAddr := addr - HRAM_START
	return func() *uint8 {
		return &m.hRAM[newAddr]
	}
}

func (m *memory) read_IE() readMemFunc {
	return func() *uint8 {
		return &m.IE[0]
	}
}
