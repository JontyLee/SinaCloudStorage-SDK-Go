/**
 * Bucket or bucket/object to delete
 * File  : cli/delete.go
 * Author: JontyLee
 * Date  : 2025-02-17 19:32:50
 */
package main

import "github.com/urfave/cli/v2"

func delete(cliCtx *cli.Context) error {
	obj := cliCtx.String("object")
	if obj == "" {
		return retry(cliCtx, func(ctx *cli.Context) error {
			return bucketInstance.DelBucket()
		})
	}
	return retry(cliCtx, func(ctx *cli.Context) error {
		return bucketInstance.Del(obj)
	})
}
