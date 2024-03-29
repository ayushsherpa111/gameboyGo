package memory

import (
	"fmt"

	"github.com/ayushsherpa111/gameboyEMU/types"
)

// TODO: handle using `interface.cart`. Placeholder for now.
func (m *memory) write_rom(addr uint16) types.WriteMemFunc {
	return func(v uint8) error {
		m.cart.WriteROM(addr, v)
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

func (m *memory) write_oam(addr uint16, cycleCount uint64) types.WriteMemFunc {
	newAddr := addr - OAM_START

	return func(val uint8) error {
		m.gpu.Write_OAM(newAddr, val, true)
		return nil
	}
}

func (m *memory) ignore_io_write() types.WriteMemFunc {
	return func(u uint8) error {
		return nil
	}
}

func (m *memory) handleDMA(cycleCount uint64) types.WriteMemFunc {
	return func(u uint8) error {
		var address uint16 = uint16(u) << 8
		var i uint16
		for i = 0; i < 0xA0; i++ {
			m.gpu.Write_OAM(i, *m.getReadMemBlock(address+i, cycleCount)(), true)
		}
		return nil
	}
}

func (m *memory) write_io(addr uint16, cycleCount uint64) types.WriteMemFunc {
	if _, ok := PPU_REGS[addr]; ok {
		return func(val uint8) error {
			return m.gpu.Write_Regs(addr, val)
		}
	}

	if addr == DMA {
		return m.handleDMA(cycleCount)
	}

	newAddr := addr - IO_START
	return func(u uint8) error {
		m.ioRegs[newAddr] = u
		if newAddr == 0x50 {
			m.isBootLoaderLoaded = false
		}

		if newAddr == 0x00 {
			m.joypadCtx.SetSelBit(u)
		}

		if newAddr == 0x04 {
			var DIV uint8 = m.ioRegs[newAddr]
			m.ioRegs[newAddr] = DIV + uint8(((uint64(DIV)*256)-cycleCount)/256)
		}

		switch newAddr {
		case TIMA - IO_START:
			fmt.Printf("TIMA : 0x%x CYCLE COUNT %d\n", u, cycleCount)
			m.lastCycleCount = cycleCount

			m.Scheduler.ClearEventType(types.EV_TIMER)
			m.scheduleTimerEvents(u, cycleCount)
		case TAC - IO_START:
			fmt.Printf("TAC: 0x%x\n", u)
			m.ioRegs[TIMA-IO_START] += uint8((cycleCount - m.lastCycleCount) / m.getClockTiming())
			m.lastCycleCount = cycleCount

			m.Scheduler.ClearEventType(types.EV_TIMER)
			m.scheduleTimerEvents(m.ioRegs[TIMA-IO_START], cycleCount)
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
