package memory

import (
	"github.com/ayushsherpa111/gameboyEMU/types"
)

// TODO: handle using `interface.cart`. Placeholder for now.
func (m *memory) write_rom(addr uint16) types.WriteMemFunc {
	return func(val uint8) error {
		return nil
	}
}

func (m *memory) write_vram(addr uint16) types.WriteMemFunc {
	newAddr := addr - VRAM_START
	return func(val uint8) error {
		m.gpu.Write_VRAM(newAddr, val)
		return nil
	}
}

// TODO: implement through cart.
func (m *memory) write_eram(addr uint16) types.WriteMemFunc {
	newAddr := addr - EXT_RAM_START
	return func(val uint8) error {
		m.eRAM[newAddr] = val
		return nil
	}
}

func (m *memory) write_wram(addr uint16) types.WriteMemFunc {
	newAddr := mapwRAMIndex(addr)
	return func(u uint8) error {
		m.wRAM[newAddr] = u
		return nil
	}
}

func (m *memory) write_oam(addr uint16) types.WriteMemFunc {
	newAddr := addr - OAM_START
	return func(val uint8) error {
		m.gpu.Write_OAM(newAddr, val)
		return nil
	}
}

func (m *memory) ignore_io_write() types.WriteMemFunc {
	return func(u uint8) error {
		return nil
	}
}

func (m *memory) write_io(addr uint16, cycleCount uint64) types.WriteMemFunc {
	if _, ok := PPU_REGS[addr]; ok {
		return func(val uint8) error {
			return m.gpu.Write_Regs(addr, val)
		}
	}

	newAddr := addr - IO_START
	return func(u uint8) error {
		m.ioRegs[newAddr] = u
		switch newAddr {
		case TIMA - IO_START:
			m.lastCycleCount = cycleCount

			m.Scheduler.ClearEventType(types.EV_TIMER)
			m.scheduleTimerEvents(u)
		case TAC - IO_START:
			m.lastCycleCount = cycleCount

			m.Scheduler.ClearEventType(types.EV_TIMER)
			m.scheduleTimerEvents(m.ioRegs[TIMA-IO_START])
		}

		return nil
	}
}

func (m *memory) write_hram(addr uint16) types.WriteMemFunc {
	newAddr := addr - HRAM_START
	return func(u uint8) error {
		m.hRAM[newAddr] = u
		return nil
	}
}

func (m *memory) write_IE(addr uint16) types.WriteMemFunc {
	return func(u uint8) error {
		m.IE[0] = u
		return nil
	}
}

func (m *memory) write_IF(addr uint16) types.WriteMemFunc {
	return func(u uint8) error {
		m.IF[0] = u
		return nil
	}
}
