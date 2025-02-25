/**
 * 获取文件列表
 * File  : cli/cmd/list.go
 * Author: JontyLee
 * Date  : 2025-02-24 14:26:30
 */
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists owned buckets or list bucket contents",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var bucket string
		if len(args) > 0 {
			bucket = args[0]
		}
		bucketInstance = s3.Bucket(bucket)
		if bucket == "" {
			return retry(func() error {
				data, err := bucketInstance.ListBucket()
				if err != nil {
					return err
				}
				fmt.Fprintln(os.Stdout, string(data))
				return nil
			})
		}
		return retry(func() error {
			data, err := bucketInstance.ListObject(prefix, delimiter, marker, maxKeys)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		})
	},
}

func init() {
	rootCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "prefix for results set")
	rootCmd.Flags().StringVarP(&marker, "marker", "k", "", "where in results set to start listing")
	rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", "", "delimiter for rolling up results set")
	rootCmd.Flags().IntVarP(&maxKeys, "maxKeys", "m", 0, "maximum number of keys to return in results set")

	rootCmd.AddCommand(listCmd)
}
