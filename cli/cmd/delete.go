/**
 * 删除bucket或者bucket下的对象
 * File  : cli/cmd/delete.go
 * Author: JontyLee
 * Date  : 2025-02-25 11:37:35
 */
package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a bucket or bucket/object",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		bucketInstance = s3.Bucket(bucket)
		if len(args) == 1 {
			return retry(func() error {
				return bucketInstance.DelBucket()
			})
		}
		return retry(func() error {
			return bucketInstance.Del(args[1])
		})
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
