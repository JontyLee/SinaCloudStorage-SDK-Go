/**
 * 创建新的bucket
 * File  : cli/cmd/create.go
 * Author: JontyLee
 * Date  : 2025-02-24 17:19:02
 */
package cmd

import (
	"github.com/spf13/cobra"

	sinastoragegosdk "github.com/SinaCloudStorage/SinaCloudStorage-SDK-Go"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new bucket",
	Args:    cobra.ExactArgs(1),
	PreRunE: checkAcl,
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		bucketInstance = s3.Bucket(bucket)
		return retry(func() error {
			return bucketInstance.PutBucket(acl)
		})
	},
}

func init() {
	createCmd.Flags().StringVarP(&cannedAcl, "cannedAcl", "a", string(sinastoragegosdk.Private), "canned ACL for the bucket (see Canned ACLs)")
	rootCmd.AddCommand(createCmd)
}
