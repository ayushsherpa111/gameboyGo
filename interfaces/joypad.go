package interfaces

type Joypad interface {
	/*
		Set bits 5 or 4 to denote the usage of either action or direction buttons
	*/
	SetSelBit(uint8)

	GetGamepadState() *uint8
}
