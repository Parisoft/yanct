package cmd

import (
	"errors"

	"github.com/parisoft/yanct/chr"

	"github.com/spf13/cobra"
)

//Metasrpite output format
const (
	MetaspriteOutputC   = "c"
	MetaspriteOutputASM = "asm"
	MetaspriteOutputBin = "bin"
)

var img2sprCmd = &cobra.Command{
	Use:   "img2spr IMAGE",
	Short: "Convert a PNG image into a CHR + Metasprite file",
	Long: `Convert a PNG image into a CHR + Metasprite file.
First the image is converted into a CHR containing tiles of the choosen dimension, then all blank and duplicated tiles are removed.
A metasprite file is also generated into the choosen format with the (0,0) axis pointing to the bottom left corner of the image.
The image must be indexed with 4 colors and has the maximum dimension of 128x128 pixels.`,
	Example: `Convert the image 'sprite.png' into a CHR with 8x16 tiles and a metasprite formatted as C source code.
This command will generate 1 file for CHR: sprite.chr and 2 files for metasprite: sprite.c and sprite.h

yanct im2spr sprite.png --tile-height=16 --metasprite-format=c`,
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
		if err := validateMetasprFmt(); err != nil {
			return err
		}
		if err := validateBgColor(); err != nil {
			return err
		}
		if err := validateTileH(); err != nil {
			return err
		}
		return convert(args[0])
	},
}

func init() {
	img2sprCmd.Flags().Uint8VarP(&flg.bgColor, FlgBgColor, "b", 0, "Color index of the background [0,3] (default 0)")
	img2sprCmd.Flags().Uint8VarP(&flg.tileH, FlgTileH, "t", 8, "Height of the tiles: 8 for 8x8, 16 for 8x16")
	img2sprCmd.Flags().StringVarP(&flg.metasprFmt, FlgMetasprFmt, "f", "bin", "Metasprite output format: c, asm, bin")
	rootCmd.AddCommand(img2sprCmd)
}

func convert(filename string) error {
	pngimg, err := openImg(filename)
	if err != nil {
		return err
	}

	tileset := chr.NewTilesetFromPNG(pngimg, flg.bgColor)
	metasprite := chr.NewMetaspriteFromTileset(tileset)

	if flg.tileH == 16 {
		tileset.To8x16()
		metasprite.To8x16()
	}

	chr.CleanupTiles(tileset, metasprite)

	err = tileset.Write(filename)
	if err != nil {
		return err
	}

	switch flg.metasprFmt {
	case MetaspriteOutputASM:
		err = metasprite.WriteAsm(filename)
	case MetaspriteOutputBin:
		err = metasprite.WriteBin(filename)
	case MetaspriteOutputC:
		err = metasprite.WriteC(filename)
	}

	return err
}
