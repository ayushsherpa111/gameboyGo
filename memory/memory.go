package memory

type Memory struct {
	memory [0xFFFF]uint8
}

func InitMem() *Memory {
	return &Memory{
		memory: [0xFFFF]uint8{},
	}
}

func (m *Memory) GetByte(addr uint16) byte {
	return m.memory[addr]
}

func (m *Memory) SetByte(addr uint16, val byte) {
	m.memory[addr] = val
}
