/**
 * 复制1个对象
 * File  : cli/cmd/copy.go
 * Author: JontyLee
 * Date  : 2025-02-25 12:00:51
 */
package cmd

import "github.com/spf13/cobra"

// Todo:
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "copies an object; if any options are set, the entire metadata of the object is replaced",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	// rootCmd.AddCommand(copyCmd)
}
