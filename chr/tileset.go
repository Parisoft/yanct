package chr

import (
	"image"
	"os"
)

const (
	//TilesetMaxRows is the max number of rows
	TilesetMaxRows = 16
	//TilesetMaxCols is the max number of cols
	TilesetMaxCols = 16
)

//Tileset is a table of tiles
type Tileset struct {
	tiles   []*Tile
	tiledim TileDimension
}

//NewTilesetFromPNG builds a tileset from an indexed PNG image
func NewTilesetFromPNG(img image.PalettedImage, bgColorIdx byte) *Tileset {
	tileset := new(Tileset)
	tileset.tiledim = Tile8x8
	for i := 0; i < 256; i++ {
		tileset.tiles = append(tileset.tiles, &Tile{})
	}
	h := byte(img.Bounds().Dy())
	w := byte(img.Bounds().Dx())

	for row, y := TilesetMaxRows-h/8, byte(0); y < h; row, y = row+1, y+8 {
		for col, x := byte(0), byte(0); x < w; col, x = col+1, x+8 {
			tile := tileset.At(int(row*TilesetMaxCols + col))
			pixels := pixels(x, y, img)

			//http://wiki.nesdev.com/w/index.php/PPU_pattern_tables
			for i := byte(0); i < 8; i++ {
				for j := byte(0); j < 8; j++ {
					pixel := pixels[i*8+j]
					if bgColorIdx > 0 {
						if pixel == byte(bgColorIdx) {
							pixel = 0
						} else if pixel == 0 {
							pixel = byte(bgColorIdx)
						}
					}

					tile.Plane[0][i] |= (pixel & 1) << (7 - j)
					tile.Plane[1][i] |= ((pixel & 2) >> 1) << (7 - j)
				}
			}
		}
	}

	return tileset
}

//NewTilesetFromCHR builds a tileset from a CHR file
func NewTilesetFromCHR(chrfile *os.File, dim TileDimension) (*Tileset, error) {
	stat, err := chrfile.Stat()
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, stat.Size())
	if _, err := chrfile.Read(bytes); err != nil {
		return nil, err
	}

	tileset := NewTileset(dim)
	for i := 0; i < len(bytes); i += 16 {
		tile := new(Tile)
		tileset.tiles = append(tileset.tiles, tile)
		for b := 0; b < 8; b++ {
			tile.Plane[0][b] = bytes[i+b]
			tile.Plane[1][b] = bytes[i+b+8]
		}
	}

	return tileset, nil
}

//NewTileset builds an empty tileset for a given tile dimension
func NewTileset(tiledim TileDimension) *Tileset {
	return &Tileset{tiledim: tiledim}
}

//To8x16 convert the tiles to 8x16 pixels
func (tileset *Tileset) To8x16() {
	tmp := make([]*Tile, len(tileset.tiles))

	// move sprites at odd lines to 1 line above at odd column
	for row, i := 0, -1; row < tileset.Size()/TilesetMaxRows; row += 2 {
		for col := 0; col < TilesetMaxCols; col++ {
			idx := row*TilesetMaxCols + col
			i++
			tmp[i] = tileset.At(idx)
			i++
			tmp[i] = tileset.At(idx + TilesetMaxCols)
		}
	}

	tileset.tiles = tmp
	tileset.tiledim = Tile8x16
}

func (tileset Tileset) Write(filename string) error {
	chrfile := changeFileExtension(filename, "chr")
	file, err := os.OpenFile(chrfile, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	for _, tile := range tileset.tiles {
		if _, err := file.Write(tile.Plane[0][:]); err != nil {
			return err
		}
		if _, err := file.Write(tile.Plane[1][:]); err != nil {
			return err
		}
	}

	return nil
}

//At returns a tile from tileset at position i
func (tileset *Tileset) At(i int) *Tile {
	return tileset.tiles[i]
}

//RemoveAt remove an tile at position i
func (tileset *Tileset) RemoveAt(i int) {
	tileset.tiles = append(tileset.tiles[:i], tileset.tiles[i+1:]...)
}

//Size returns how many tiles the tileset contains
func (tileset *Tileset) Size() int {
	return len(tileset.tiles)
}

func pixels(x, y byte, img image.PalettedImage) []byte {
	pixels := make([]byte, 64)

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			pixels[i*8+j] = img.ColorIndexAt(int(x)+j, int(y)+i)
		}
	}

	return pixels
}
