package cmd

import (
	"errors"
	"os"

	"github.com/parisoft/yanct/chr"

	"github.com/spf13/cobra"
)

var concatCmd = &cobra.Command{
	Use:   "concat CHR_1 CHR_2 [...CHR_N]",
	Short: "Concatenate many CHR files into one",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("concat requires 2 CHR files or more")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateMetasprFmt(); err != nil {
			return err
		}
		if err := validateTileH(); err != nil {
			return err
		}
		if err := validateChrOut(); err != nil {
			return err
		}
		
		return concat(args...)
	},
}

func init() {
	concatCmd.Flags().Uint8VarP(&flg.tileH, FlgTileH, "t", 8, "Height of the tiles: 8 for 8x8, 16 for 8x16")
	concatCmd.Flags().StringVarP(&flg.chrOut, FlgChrOut, "o", "", "output CHR file name")
	concatCmd.MarkFlagRequired(FlgChrOut)
	rootCmd.AddCommand(concatCmd)
}

func concat(chrlist ...string) error {
	tilesets := make([]*chr.Tileset, len(chrlist))
	metasprites := make([]*chr.Metasprite, len(chrlist))
	binnames := make([]string, len(chrlist))
	for _, chrfilename := range chrlist {
		chrfile, err := os.Open(chrfilename)
		if err != nil {
			return err
		}
		defer chrfile.Close()

		binfilename := chrfilename[:len(chrfilename)-3] + "bin"
		binnames = append(binnames, binfilename)
		binfile, err := os.Open(binfilename)
		if err == nil {
			defer binfile.Close()
			metasprites = append(metasprites, chr.NewMetaspriteFromFile(binfile))
		} else if os.IsNotExist(err) {
			metasprites = append(metasprites, nil)
		} else {
			return err
		}

		tilesets = append(tilesets, chr.NewTilesetFromCHR(chrfile))
	}

	pivot := new(chr.Tileset)
	if flg.tileH == 16 {
		pivot.To8x16()

		for i, tileset := range tilesets {
			tileset.To8x16()

			if metasprite := metasprites[i]; metasprite != nil {
				metasprite.To8x16()
			}
		}
	}

	for i, tileset := range tilesets {
		chr.ConcatTiles(pivot, tileset, metasprites[i])
	}

	pivot.Write(flg.chrOut)

	for i, metasprite := range metasprites {
		if metasprite != nil {
			metasprite.WriteBin(binnames[i])
		}
	}

	return nil
}
