/**
 * 获取1个对象的头信息
 * File  : cli/cmd/head.go
 * Author: JontyLee
 * Date  : 2025-02-25 12:14:23
 */
package cmd

import (
	"github.com/spf13/cobra"
)

// Todo:
var headCmd = &cobra.Command{
	Use:   "gqs",
	Short: "generates an authenticated query string",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	// rootCmd.AddCommand(headCmd)
}
