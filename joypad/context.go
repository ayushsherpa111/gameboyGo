package joypad

const (
	actionBit    uint8 = 0x20
	directionBit uint8 = 0x10
)

type context struct {

	/*
		Selection bit
		Bit 5. (0xFF00)
	*/
	actionBit bool

	/*
		Selection bit
		Bit 4. (0xFF00)
	*/
	directionBit bool

	controller joypad
}

func (c context) SetSelBit(val uint8) {
	if (^val & actionBit) > 0 {
		c.actionBit = true
	} else {
		c.actionBit = false
	}

	if (^val & directionBit) > 0 {
		c.directionBit = true
	} else {
		c.directionBit = false
	}
}

func (c context) GetGamepadState() *uint8 {
	if c.directionBit {
		return &c.controller.DirectionBits
	}
	if c.actionBit {
		return &c.controller.ActionBits
	}
	return &c.controller.Default
}

func NewContext() context {
	return context{
		actionBit:    false,
		directionBit: false,
		controller: joypad{
			DirectionBits: 0x0F,
			ActionBits:    0x0F,
			Default:       0xFF,
		},
	}
}

func (c context) HandleEvent(key uint8, state bool) {
	if (key & directionBit) != 0 {
		c.controller.SetDirection(key&0xF, state)
	} else if (key & actionBit) != 0 {
		c.controller.SetAction(key&0xF, state)
	}
}
