package mbc

type mbc1 struct {
	bank      []uint8
	bankCount uint8
}

func (m *mbc1) GetSlice(bankNum uint8) []uint8 {
	var endIDX uint
	var startIDX uint = uint(bankNum-1) * (16 << 10)
	if bankNum == m.bankCount {
		endIDX = uint(len(m.bank))
	} else {
		endIDX = startIDX + (16 << 10)
	}
	return m.bank[startIDX:endIDX]
}
