/**
 * Get the ACL of a bucket or a bucket/object
 * File  : cli/getacl.go
 * Author: JontyLee
 * Date  : 2025-02-17 19:59:06
 */
package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func getacl(cliCtx *cli.Context) error {
	if hasWriter {
		defer output.Close()
	}
	if object == "" {
		return retry(cliCtx, func(ctx *cli.Context) error {
			data, err := bucketInstance.GetBucketInfo("acl")
			if err != nil {
				return err
			}
			fmt.Fprintf(output, "%s\n", data)
			return nil
		})
	}
	return retry(cliCtx, func(ctx *cli.Context) error {
		data, err := bucketInstance.GetInfo(object, "acl")
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "%s\n", data)
		return nil
	})
}
