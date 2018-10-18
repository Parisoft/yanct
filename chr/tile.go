package chr

//Tile define a 8x8 pixels of a CHR file
type Tile struct {
	Plane [2][8]byte
}

//TileDimension is the dimension in pixels of a tile
type TileDimension string

const(
	//Tile8x8 have 8x8 pixels
	Tile8x8 TileDimension = "8x8"
	//Tile8x16 have 8x16 pixels
	Tile8x16 TileDimension = "8x16"
)

//NewTile constructs a new Tile
func NewTile(bytes []byte) Tile {
	tile := Tile{}
	copy(tile.Plane[0][:], bytes[0:8])
	copy(tile.Plane[1][:], bytes[8:16])

	return tile
}

//Equals returns true if 2 tiles has the same planes
func (tile *Tile) Equals(other *Tile) bool {
	for p := 0; p < 2; p++ {
		for b := 0; b < 8; b++ {
			if tile.Plane[p][b] != other.Plane[p][b] {
				return false
			}
		}
	}

	return true
}

//Empty returns true if the tile planes contains only zeroes
func (tile *Tile) Empty() bool {
	for _, plane := range tile.Plane {
		for _, b := range plane {
			if b != 0 {
				return false
			}
		}
	}

	return true
}
