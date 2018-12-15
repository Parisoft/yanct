package cmd

import (
	"errors"
	"os"

	"github.com/parisoft/yanct/chr"
	"github.com/spf13/cobra"
)

var mergemetaCmd = &cobra.Command{
	Use:   "mergemeta METASPR_1 METASPR_2 [METASPR_N]",
	Short: "Merge many metasprite files into one",
	Long: `Concatenate many metasprite files into one.
All files are appended to the first one.
Only the binary format is allowed.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("mergemeta requires 2 metasprite files or more")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOutFileName(); err != nil {
			return err
		}

		return mergemeta(args...)
	},
}

func init() {
	mergemetaCmd.Flags().StringVarP(&flg.fileOut, FlgOutFile, "o", "", "output metasprite file name")
	rootCmd.AddCommand(mergemetaCmd)
}

func mergemeta(filenames ...string) error {
	metasprites := make([]*chr.Metasprite, len(filenames))
	for i, filename := range filenames {
		binfile, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer binfile.Close()

		metasprites[i], err = chr.NewMetaspriteFromFile(binfile)
		if err != nil {
			return err
		}
	}

	for i := 1; i < len(metasprites); i++ {
		metasprites[0].Merge(metasprites[i])
	}

	output := flg.fileOut
	if len(output) == 0 {
		output = filenames[0]
	}
	metasprites[0].WriteBin(output)

	return nil
}
