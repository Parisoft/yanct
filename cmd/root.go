package cmd

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/spf13/cobra"
)

//Flag names
const (
	FlgPal        = "pal"
	FlgBgColor    = "bg-color"
	FlgTileH      = "tile-height"
	FlgMetasprFmt = "metasprite-format"
	FlgChrOut     = "chr-output"
)

type flag struct {
	pal        uint8
	bgColor    uint8
	tileH      uint8
	metasprFmt string
	chrOut     string
}

var flg flag

var rootCmd = &cobra.Command{
	Use:   "yanct",
	Short: "Yet Another NES CHR Tool",
	Long:  "A command line tool to create and edit CHR files",
}

//Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validatePal() error {
	if flg.pal > 3 {
		return fmt.Errorf("Invalid palette index (%s): %d", FlgPal, flg.pal)
	}
	return nil
}

func validateBgColor() error {
	if flg.bgColor > 3 {
		return fmt.Errorf("Invalid background color index (%s): %d", FlgBgColor, flg.bgColor)
	}
	return nil
}

func validateTileH() error {
	if flg.tileH != 8 && flg.tileH != 16 {
		return fmt.Errorf("Invalid tile height (%s): %d", FlgTileH, flg.tileH)
	}
	return nil
}

func validateMetasprFmt() error {
	if flg.metasprFmt != MetaspriteOutputC && flg.metasprFmt != MetaspriteOutputASM && flg.metasprFmt != MetaspriteOutputBin {
		return fmt.Errorf("Invalid metasprite output format (%s): %s", FlgMetasprFmt, flg.metasprFmt)
	}
	return nil
}

func validateChrOut() error {
	if len(flg.chrOut) == 0 {
		return fmt.Errorf("Invalid output CHR file name (%s): %s", FlgChrOut, flg.chrOut)
	}
	return nil
}

func openImg(filename string) (image.PalettedImage, error) {
	pngfile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer pngfile.Close()

	decoded, err := png.Decode(pngfile)
	if err != nil {
		return nil, err
	}

	img, ok := decoded.(image.PalettedImage)
	if !ok || img.Bounds().Dx() > 128 || img.Bounds().Dy() > 128 {
		return nil, fmt.Errorf("Image '%s' must be a PNG file indexed with 4 colors and has the maximum dimension of 128x128 pixels", filename)
	}

	return img, nil
}
