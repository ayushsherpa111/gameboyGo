package mbc

type mbc2 struct {
	bank      []uint8
	bankCount uint8
}

func (m *mbc2) GetSlice(_ uint8) []uint8 {
	panic("not implemented") // TODO: Implement
}
