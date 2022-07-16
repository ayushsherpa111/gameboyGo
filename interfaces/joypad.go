package interfaces

type Joypad interface {
	/*
		Set bits 5 or 4 to denote the usage of either action or direction buttons
	*/
	SetSelBit(uint8)

	GetGamepadState() *uint8

	/*
		Register Bits related to Direction/Action as key press (0 = Pressed)
	*/
	KeyDown(uint8)

	/*
		Register Bits related to Direction/Action as key release (0 = Pressed)
	*/
	KeyUp(uint8)
}
