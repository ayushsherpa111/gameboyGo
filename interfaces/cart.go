package interfaces

type Cart interface {
	HeaderInfo()
	ReadROM(uint16) *uint8
	WriteROM(uint16, uint8)
}
