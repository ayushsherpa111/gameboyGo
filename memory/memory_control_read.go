package memory

import "github.com/ayushsherpa111/gameboyEMU/types"

func (m *memory) ignore_io_read() types.ReadMemFunc {
	return func() *uint8 { return nil }
}

func (m *memory) read_io(addr uint16, cycleCount uint64) types.ReadMemFunc {
	if _, ok := PPU_REGS[addr]; ok {
		return func() *uint8 {
			return m.gpu.Read_Regs(addr)
		}
	}

	newAddr := addr - IO_START
	return func() *uint8 {

		if newAddr == 0x00 {
			return m.joypadCtx.GetGamepadState()
		}

		switch addr {
		case TIMA:
			// should retrn a value between 0x0 - 0xFF
			m.ioRegs[TIMA-IO_START] += uint8((cycleCount - m.lastCycleCount) / m.getClockTiming())
		}
		return &m.ioRegs[newAddr]
	}
}

func (m *memory) read_boot_loader(addr uint16) types.ReadMemFunc {
	return func() *uint8 {
		return &m.bootloader[addr]
	}
}

func (m *memory) read_rom_data(addr uint16) types.ReadMemFunc {
	return func() *uint8 {
		return &m.romData[addr]
	}
}

func (m *memory) read_vram_data(addr uint16) types.ReadMemFunc {
	newAddr := addr - VRAM_START
	return func() *uint8 {
		return m.gpu.Read_VRAM(newAddr)
	}
}

// TODO: EXT RAM from Cartridge
func (m *memory) read_ext_ram(addr uint16) types.ReadMemFunc {
	newAddr := addr - EXT_RAM_START
	return func() *uint8 {
		return &m.eRAM[newAddr]
	}
}

func (m *memory) read_wram(addr uint16) types.ReadMemFunc {
	newAddr := mapwRAMIndex(addr)
	return func() *uint8 {
		return &m.wRAM[newAddr]
	}
}

func (m *memory) read_oam(addr uint16) types.ReadMemFunc {
	newAddr := addr - OAM_START
	return func() *uint8 {
		return m.gpu.Read_OAM(newAddr)
	}
}

func (m *memory) read_hram(addr uint16) types.ReadMemFunc {
	newAddr := addr - HRAM_START
	return func() *uint8 {
		return &m.hRAM[newAddr]
	}
}

func (m *memory) read_IE() types.ReadMemFunc {
	return func() *uint8 {
		return &m.IE[0]
	}
}
