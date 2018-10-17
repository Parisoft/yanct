package chr

import (
	"fmt"
)

// Sprite defines a sprite
type Sprite struct {
	X   int8
	Y   int8
	Opt byte
	Idx byte
}

// Bytes transform a Sprite into an array of bytes in the format [x, y, idx, opt]
func (spr *Sprite) Bytes() []byte {
	bytes := make([]byte, 4)
	bytes[0] = byte(spr.X)
	bytes[1] = byte(spr.Y)
	bytes[2] = spr.Idx
	bytes[3] = spr.Opt

	return bytes
}

func (spr *Sprite) String() string {
	return fmt.Sprintf("%d, %d, 0x%x, %d", spr.X, spr.Y, spr.Idx, spr.Opt)
}
