package memory

import "errors"

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

func (m *Memory) GetByte(addr uint16) (error, byte) {
	if err := validateMemAddr(int(addr), len(m.memory)); err != nil {
		return err, 0
	}
	return nil, m.memory[addr]
}

func (m *Memory) SetByte(addr uint16, val byte) error {
	if err := validateMemAddr(int(addr), len(m.memory)); err != nil {
		return err
	}
	m.memory[addr] = val
	return nil
}
