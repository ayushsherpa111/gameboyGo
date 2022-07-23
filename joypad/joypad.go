package joypad

const (
	joypadMask  uint8 = 0x0F
	downOrStart uint8 = 0x08
	upOrSelect  uint8 = 0x04
	leftOrB     uint8 = 0x02
	rightOrA    uint8 = 0x01
)

type joypad struct {
	Default uint8
	/*
		Get Directions when bit 4 is set
		0 = pressed
	*/
	DirectionBits uint8

	/*
		Get Actions when bit 5 is set
		0 = pressed
	*/
	ActionBits uint8
}

func unmask(val, mask uint8) bool {
	if ^val&mask == mask {
		return true
	}
	return false
}

func (j joypad) SetDirection(direction uint8, state bool) {
	if state {
		j.DirectionBits &= direction
	} else {
		j.DirectionBits |= ((^direction) & 0x0F)
	}
}

func (j joypad) SetAction(action uint8, state bool) {
	if state {
		j.ActionBits &= action
	} else {
		j.ActionBits |= ((^action) & 0x0F)
	}
}
