package mbc

type mbc3 struct {
	bank      []uint8
	bankCount uint8
}

func (m *mbc3) GetSlice(_ uint8) []uint8 {
	return m.bank[:]
}
