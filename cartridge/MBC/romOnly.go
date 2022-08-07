package mbc

type romOnly struct {
	bank []uint8
}

func (m *romOnly) GetSlice(_ uint8) []uint8 {
	return m.bank[:]
}
