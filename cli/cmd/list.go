/**
 * 获取文件列表
 * File  : cli/cmd/list.go
 * Author: jianlin6
 * Date  : 2025-02-24 14:26:30
 */
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists owned buckets or list bucket contents",
	RunE: func(cmd *cobra.Command, args []string) error {
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
		options := splitArgs(args, 1)
		prefix := options["prefix"]
		delimiter := options["delimiter"]
		marker := options["marker"]
		var maxKeys int
		if options["maxKeys"] != "" {
			var errMaxKeys error
			maxKeys, errMaxKeys = strconv.Atoi(options["maxKeys"])
			if errMaxKeys != nil {
				return errMaxKeys
			}
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
	rootCmd.AddCommand(listCmd)
}
