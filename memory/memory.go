package memory

import "errors"

// 0x0000 - 0x3FFF : ROM Bank 0
// 0x4000 - 0x7FFF : ROM Bank 1 - Switchable
// 0x8000 - 0x97FF : CHR RAM
// 0x9800 - 0x9BFF : BG Map 1
// 0x9C00 - 0x9FFF : BG Map 2
// 0xA000 - 0xBFFF : Cartridge RAM
// 0xC000 - 0xCFFF : RAM Bank 0
// 0xD000 - 0xDFFF : RAM Bank 1-7 - switchable - Color only
// 0xE000 - 0xFDFF : Reserved - Echo RAM
// 0xFE00 - 0xFE9F : Object Attribute Memory
// 0xFEA0 - 0xFEFF : Reserved - Unusable
// 0xFF00 - 0xFF7F : I/O Registers
// 0xFF80 - 0xFFFE : Zero Page

type Memory struct {
	memory [0xFFFF]uint8
}

func InitMem() *Memory {
	return &Memory{
		memory: [0xFFFF]uint8{},
	}
}

func validateMemAddr(addr, memLen int) error {
	if addr > memLen || addr < 0 {
		return errors.New("Invalid Memory Range")
	}
	return nil
}

func (m *Memory) GetByte(addr uint16) (error, *byte) {
	if err := validateMemAddr(int(addr), len(m.memory)); err != nil {
		return err, nil
	}
	return nil, &m.memory[addr]
}

func (m *Memory) SetByte(addr uint16, val byte) error {
	if err := validateMemAddr(int(addr), len(m.memory)); err != nil {
		return err
	}
	m.memory[addr] = val
	return nil
}
