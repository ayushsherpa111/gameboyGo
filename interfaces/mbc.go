package interfaces

type MBC interface {
	// Get slice returns a slice of 16KB from
	// its internal bank depending on the bank number argument passed
	GetSlice(uint8) []uint8
}
