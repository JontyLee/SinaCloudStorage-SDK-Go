/**
 * 获取1个对象
 * File  : cli/cmd/get.go
 * Author: JontyLee
 * Date  : 2025-02-25 11:54:00
 */
package cmd

import (
	"github.com/spf13/cobra"
)

// Todo:
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "gets an object",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// bucket := args[0]
		// bucketInstance = s3.Bucket(bucket)
		// if len(args) == 1 {
		// 	return retry(func() error {
		// 		return bucketInstance.Get(args[1])
		// 	})
		// }
		// return retry(func() error {
		// 	return bucketInstance.Del(args[1])
		// })
		return nil
	},
}

func init() {
	// rootCmd.AddCommand(getCmd)
}
