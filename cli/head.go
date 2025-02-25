/**
 * Gets only the headers of a bucket or bucket object
 * File  : cli/head.go
 * Author: jianlin6
 * Date  : 2025-02-18 11:25:40
 */
package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func head(cliCtx *cli.Context) error {
	obj := cliCtx.String("object")
	if obj == "" {
		return retry(cliCtx, func(ctx *cli.Context) error {
			data, err := bucketInstance.GetBucketInfo("meta")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", data)
			return nil
		})
	}
	return retry(cliCtx, func(ctx *cli.Context) error {
		data, err := bucketInstance.GetInfo(obj, "meta")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", data)
		return nil
	})
}
