/**
 * Bucket or bucket/object to delete
 * File  : cli/delete.go
 * Author: JontyLee
 * Date  : 2025-02-17 19:32:50
 */
package cmd

import "github.com/urfave/cli/v2"

func delete(cliCtx *cli.Context) error {
	if object == "" {
		return retry(cliCtx, func(ctx *cli.Context) error {
			return bucketInstance.DelBucket()
		})
	}
	return retry(cliCtx, func(ctx *cli.Context) error {
		return bucketInstance.Del(object)
	})
}
