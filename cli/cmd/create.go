/**
 * 创建新的bucket
 * File  : cli/cmd/create.go
 * Author: jianlin6
 * Date  : 2025-02-24 17:19:02
 */
package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket = args[0]
		bucketInstance = s3.Bucket(bucket)
		options := splitArgs(args, 1)
		if err := validateAcl(options); err != nil {
			return err
		}
		return retry(func() error {
			return bucketInstance.PutBucket(acl)
		})
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
