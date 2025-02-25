/**
 * 获取指定bucket或者对象的权限设置信息
 * File  : cli/cmd/getacl.go
 * Author: JontyLee
 * Date  : 2025-02-25 12:12:30
 */
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var getaclCmd = &cobra.Command{
	Use:   "getacl",
	Short: "get the ACL of a bucket or object",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		bucketInstance = s3.Bucket(bucket)
		if len(args) == 1 {
			return retry(func() error {
				data, err := bucketInstance.GetBucketInfo("acl")
				if err != nil {
					return err
				}
				_, err = fmt.Fprintln(os.Stdout, string(data))
				if err != nil {
					return err
				}
				return nil
			})
		}
		return retry(func() error {
			data, err := bucketInstance.GetInfo(args[1], "acl")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(os.Stdout, string(data))
			if err != nil {
				return err
			}
			return nil
		})
	},
}

func init() {
	// Todo:
	getaclCmd.Flags().StringVarP(&filename, "filename", "f", "", "output filename for ACL (default is stdout)")
	rootCmd.AddCommand(getaclCmd)
}
