package cmd

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/parisoft/yanct/chr"

	"github.com/spf13/cobra"
)

//Metasrpite output format
const (
	MetaspriteOutputC   = "c"
	MetaspriteOutputASM = "asm"
	MetaspriteOutputBin = "bin"
)

//Flag names
const(
	FlgBgColor ="bg-color"
	FlgTileH = "tile-height"
	FlgMetasprFmt = "metasprite-format"
)

var convertCmd = &cobra.Command{
	Use:   "convert IMAGE",
	Short: "Convert a PNG image into a CHR + Metasprite file",
	Long: `Convert a PNG image into a CHR + Metasprite file.
First the image is converted into a CHR containing tiles of the choosen dimension, then all blank and duplicated tiles are removed.
A metasprite file is also generated into the choosen format with the (0,0) axis pointing to the bottom left corner of the image.
The image must be indexed with 4 colors and has the maximum dimension of 128x128 pixels.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Missing image file name")
		}
		if len(args) > 1 {
			return errors.New("Only 1 image can be converted per execution")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return convert(args[0])
	},
}

var (
	pngimg     image.PalettedImage
	bgColor    byte
	tileH      byte
	metasprOut string
)

func init() {
	convertCmd.Flags().Uint8VarP(&bgColor, FlgBgColor, "b", 0, "Color index of the background [0,3] (default 0)")
	convertCmd.Flags().Uint8VarP(&tileH, FlgTileH, "t", 8, "Height of the tiles: 8 for 8x8, 16 for 8x16")
	convertCmd.Flags().StringVarP(&metasprOut, FlgMetasprFmt, "f", "bin", "Metasprite output format: c, asm, bin")
	rootCmd.AddCommand(convertCmd)
}

func convert(filename string) error {
	if err := validate(filename); err != nil {
		return err
	}

	tileset := chr.NewTileset(pngimg, bgColor)
	metasprite := chr.NewMetasprite(tileset)

	if tileH == 16 {
		tileset.To8x16()
		metasprite.To8x16()
		chr.Cleanup8x16Tiles(tileset, metasprite)
	} else {
		chr.Cleanup8x8Tiles(tileset, metasprite)
	}

	tileset.Write(filename)

	var err error
	switch metasprOut {
	case MetaspriteOutputASM:
		err = metasprite.WriteAsm(filename)
	case MetaspriteOutputBin:
		err = metasprite.WriteBin(filename)
	case MetaspriteOutputC:
		err = metasprite.WriteC(filename)
	}

	return err
}

func validate(filename string) error {
	if bgColor > 3 {
		return fmt.Errorf("Invalid background color index (%s): %d", FlgBgColor, bgColor)
	}
	if tileH != 8 && tileH != 16 {
		return fmt.Errorf("Invalid tile height (%s): %d", FlgTileH, tileH)
	}
	if metasprOut != MetaspriteOutputC && metasprOut != MetaspriteOutputASM && metasprOut != MetaspriteOutputBin {
		return fmt.Errorf("Invalid metasprite output format (%s): %s", FlgMetasprFmt, metasprOut)
	}

	pngfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer pngfile.Close()

	decoded, err := png.Decode(pngfile)
	if err != nil {
		return err
	}

	img, ok := decoded.(image.PalettedImage)
	if !ok || img.Bounds().Dx() > 128 || img.Bounds().Dy() > 128 {
		return errors.New("The image must be a PNG file indexed with 4 colors and has the maximum dimension of 128x128 pixels")
	}

	pngimg = img

	return nil
}
