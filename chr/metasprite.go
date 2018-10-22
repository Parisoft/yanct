package chr

import (
	"fmt"
	"os"
	"strings"
)

//Metasprite is a table of sprites
type Metasprite struct {
	sprites []*Sprite
}

//NewMetaspriteFromTileset builds a metasprite from a Tileset
func NewMetaspriteFromTileset(tileset *Tileset) *Metasprite {
	metasprite := new(Metasprite)
	rows := tileset.Size() / TilesetMaxRows

	for row := 0; row < rows; row++ {
		for col := 0; col < TilesetMaxCols; col++ {
			metasprite.sprites = append(metasprite.sprites, &Sprite{
				X:   int8(col * 8),
				Y:   int8((row - rows) * 8),
				Opt: 0,
				Idx: byte(row*TilesetMaxCols + col),
			})
		}
	}

	return metasprite
}

//NewMetaspriteFromFile builds a metasprite from a binary file
func NewMetaspriteFromFile(binfile *os.File) (*Metasprite, error) {
	stat, err := binfile.Stat()
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, stat.Size())
	if _, err := binfile.Read(bytes); err != nil {
		return nil, err
	}

	metasrp := new(Metasprite)
	for i := 0; i < len(bytes)-1; i += 4 {
		metasrp.sprites = append(metasrp.sprites, &Sprite{
			X:   int8(bytes[i]),
			Y:   int8(bytes[i+1]),
			Idx: bytes[i+2],
			Opt: bytes[i+3],
		})
	}

	return metasrp, nil
}

//To8x16 convert the sprites to 8x16 pixels
func (metasprite *Metasprite) To8x16() {
	// remove odd lines
	for i := metasprite.Size() - 1; i >= 0; i-- {
		if (metasprite.At(i).Idx/TilesetMaxRows)%2 != 0 {
			metasprite.RemoveAt(i)
		}
	}

	for _, spr := range metasprite.sprites {
		spr.Idx = (spr.Idx/16)*16 + (spr.Idx%16)*2
	}
}

//WriteC write the metasprite to a .c and .h files
func (metasprite *Metasprite) WriteC(filename string) error {
	cfilename := changeFileExtension(filename, "c")
	cfile, err := os.OpenFile(cfilename, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer cfile.Close()

	hfilename := changeFileExtension(filename, "h")
	hfile, err := os.OpenFile(hfilename, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer hfile.Close()

	varname := filename[strings.LastIndex(filename, "/")+1 : strings.LastIndex(filename, ".")]
	varname = strings.Replace(varname, "-", "_", -1)
	fmt.Fprintf(hfile, "extern char %s[%d];\n", varname, metasprite.Size()*4+1)
	fmt.Fprintf(cfile, "const char %s[] = {\n", varname)
	for _, spr := range metasprite.sprites {
		fmt.Fprintf(cfile, "\t%s,\n", spr.String())
	}
	fmt.Fprintln(cfile, "\t0x80,")
	fmt.Fprintln(cfile, "};")

	return nil
}

//WriteAsm write the metasprite to a .inc file
func (metasprite *Metasprite) WriteAsm(filename string) error {
	asmfilename := changeFileExtension(filename, "inc")
	asmfile, err := os.OpenFile(asmfilename, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer asmfile.Close()

	varname := filename[strings.LastIndex(filename, "/")+1 : strings.LastIndex(filename, ".")]
	varname = strings.Replace(varname, "-", "_", -1)
	fmt.Fprintf(asmfile, "%s:\n", varname)
	for _, spr := range metasprite.sprites {
		fmt.Fprintf(asmfile, "\t.byte %s\n", spr.String())
	}
	fmt.Fprintln(asmfile, "\t.byte $80")

	return nil
}

//WriteBin write the metasprite to a .bin file
func (metasprite *Metasprite) WriteBin(filename string) error {
	binfilename := changeFileExtension(filename, "bin")
	binfile, err := os.OpenFile(binfilename, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer binfile.Close()

	for _, spr := range metasprite.sprites {
		if _, err := binfile.Write(spr.Bytes()); err != nil {
			return err
		}
	}
	_, err = binfile.Write([]byte{0x80})

	return err
}

//At returns a sprite at position i
func (metasprite *Metasprite) At(i int) *Sprite {
	return metasprite.sprites[i]
}

//RemoveAt remove an sprite at position i
func (metasprite *Metasprite) RemoveAt(i int) {
	metasprite.sprites = append(metasprite.sprites[:i], metasprite.sprites[i+1:]...)
}

//Size returns the amount of sprites into metasprite
func (metasprite *Metasprite) Size() int {
	return len(metasprite.sprites)
}
