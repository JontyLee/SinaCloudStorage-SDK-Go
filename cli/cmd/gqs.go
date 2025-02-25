/**
 * 生成带认证的查询
 * File  : cli/cmd/gqs.go
 * Author: JontyLee
 * Date  : 2025-02-25 11:58:10
 */
package cmd

import (
	"github.com/spf13/cobra"
)

// Todo:
var gqsCmd = &cobra.Command{
	Use:   "gqs",
	Short: "generates an authenticated query string",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	// rootCmd.AddCommand(gqsCmd)
}
