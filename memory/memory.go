package memory

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ayushsherpa111/gameboyEMU/cartridge"
	"github.com/ayushsherpa111/gameboyEMU/interfaces"
	"github.com/ayushsherpa111/gameboyEMU/logger"
	"github.com/ayushsherpa111/gameboyEMU/types"
)

var (
	PPU_REGS = map[uint16]struct{}{
		0xFF40: {},
		0xFF41: {},
		0xFF42: {},
		0xFF43: {},
		0xFF44: {},
		0xFF45: {},
		0xFF47: {},
		0xFF48: {},
		0xFF49: {},
		0xFF4A: {},
		0xFF4B: {},
		0xFF4C: {},
		0xFF4D: {},
		0xFF4E: {},
	}
)

const (
	JOYPAD_ADDR = 0xFF00
	ROM_END     = 0x7FFF

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

	NU_START = 0xFEA0
	NU_END   = 0xFEFF

	IO_START = 0xFF00
	IO_END   = 0xFF7F
	DMA      = 0xFF46

	BOOT_LOADER_FLAG = 0xFF50
	LY_REG           = 0xFF44

	INTERRUPT_FLAG   = 0xFF0F
	INTERRUPT_ENABLE = 0xFFFF

	TIMA = 0xFF05
	TMA  = 0xFF06
	TAC  = 0xFF07

	TIMER_ENABLE = 0x4
)

const (
	VBLANK_IF uint8 = 1 << iota
	LCDS_IF
	TIMER_IF
	SERIAL_IF
	JOYPAD_IF
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
	bootloader         []uint8
	rom                string
	isBootLoaderLoaded bool

	romData        []uint8 // 0x0000 - 0x8000
	eRAM           []uint8 // 0xA000 - 0xBFFF
	wRAM           []uint8 // 0xC000 - 0xDFFF
	hRAM           []uint8 // 0xFF80 - 0xFFFE
	ioRegs         []uint8 // 0xFF00 - 0xFF70
	IE             []uint8 // 0xFFFF
	lgr            logger.Logger
	cart           interfaces.Cart
	gpu            interfaces.GPU
	Scheduler      interfaces.Scheduler
	lastCycleCount uint64
	joypadCtx      interfaces.Joypad
}

func (m *memory) setIF(bit uint8) {
	m.ioRegs[INTERRUPT_FLAG-IO_START] |= bit
}

func (m *memory) SetIFTimer(sc uint64) func() {
	return func() {
		m.setIF(TIMER_IF)
		// INFO: Function is called once the timer reaches the cycle count where TIMA overflows.
		// INFO: Assign TMA to TIMA when overflow occurs

		// m.ioRegs[TIMA-IO_START] = m.ioRegs[TMA-IO_START]
		m.MemWrite(TIMA, m.ioRegs[TMA-IO_START], sc)
		// m.scheduleTimerEvents(m.ioRegs[TIMA-IO_START])
	}
}

func (m *memory) UnloadBootloader() {
	if !m.isBootLoaderLoaded {
		// m.lgr.Println("Bootloader has already been unloaded")
		return
	}
	m.gpu.PrintDetails()
	m.isBootLoaderLoaded = false
}

func (m *memory) ReadIO(addr uint16) uint8 {
	addr -= uint16(IO_START)
	return m.ioRegs[addr]
}

func InitMem(bootLoader []byte, ROM string, debug bool, gpu interfaces.GPU, joypadCtx interfaces.Joypad) (*memory, error) {
	if len(bootLoader) != 256 {
		return nil, errors.New("Invalid Bootloader provided")
	}

	romData, err := ioutil.ReadFile(ROM)
	if err != nil {
		return nil, errors.New("Failed to load ROM file")
	}

	mem := &memory{
		isBootLoaderLoaded: true,
		eRAM:               make([]uint8, 8*1024),
		wRAM:               make([]uint8, 8*1024),
		hRAM:               make([]uint8, 127),
		romData:            make([]uint8, 32*1024),
		ioRegs:             make([]uint8, 128),
		bootloader:         bootLoader,
		rom:                ROM,
		IE:                 make([]uint8, 1),
		cart:               cartridge.NewCart(romData),
		lgr:                logger.NewLogger(os.Stdout, debug, "Memory"),
		gpu:                gpu,
		joypadCtx:          joypadCtx,
	}

	mem.cart.HeaderInfo()
	mem.ioRegs[LY_REG-IO_START] = 0x90
	gpu.RefInterruptFlag(&mem.ioRegs[INTERRUPT_FLAG-IO_START])

	if e := mem.loadROM(romData); e != nil {
		return nil, e
	}

	return mem, nil
}

func (m *memory) SetScheduler(sched interfaces.Scheduler) {
	m.Scheduler = sched
}

func (m *memory) loadROM(romData []byte) error {
	m.lgr.Infof("Loading ROM file...\n")

	if len(romData) > (32 << 10) {
		return errors.New("Invalid ROM length")
	}

	copy(m.romData, romData)

	m.lgr.Infof("ROM file Succesfully loaded.\n")
	return nil
}

func (m *memory) getReadMemBlock(addr uint16, cycleCount uint64) types.ReadMemFunc {
	if m.isBootLoaderLoaded && addr <= 0xff {
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
	} else if addr <= NU_END {
		return m.ignore_io_read()
	} else if addr <= IO_END {
		return m.read_io(addr, cycleCount)
	} else if addr <= HRAM_END {
		return m.read_hram(addr)
	} else if addr == INTERRUPT_ENABLE {
		return m.read_IE()
	}
	return nil
}

func (m *memory) getWriteMemBlock(addr uint16, cycleCount uint64) types.WriteMemFunc {
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
		return m.write_oam(addr, cycleCount)
	} else if addr <= NU_END {
		return m.ignore_io_write()
	} else if addr <= IO_END {
		return m.write_io(addr, cycleCount)
	} else if addr <= HRAM_END {
		return m.write_hram(addr)
	} else if addr == INTERRUPT_ENABLE {
		return m.write_IE(addr)
	}
	return nil
}

func (m *memory) MemRead(addr uint16, cycleCount uint64) *uint8 {
	mem := m.getReadMemBlock(addr, cycleCount)
	if mem == nil {
		// handle error
		return nil
	}
	return mem()
}

func (m *memory) MemWrite(addr uint16, val uint8, cycleCount uint64) error {
	mem := m.getWriteMemBlock(addr, cycleCount)
	if mem == nil {
		return errors.New("End of memory reached")
	}
	return mem(val)
}

func (m *memory) TickAllComponents(cycleCount uint64) {
	for i := 0; i < 4; i++ {
		m.gpu.UpdateGPU()
	}
}

func (m *memory) HandleInput(keyEvent types.KeyboardEvent) {
	fmt.Println("handling input")
	m.joypadCtx.HandleEvent(types.KeyMap[keyEvent.Key], keyEvent.State)
}
