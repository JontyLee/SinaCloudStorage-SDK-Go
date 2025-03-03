/**
 * Create a bucket
 * File  : cli/create.go
 * Author: JontyLee
 * Date  : 2025-02-17 19:14:19
 */
package cmd

import (
	"github.com/urfave/cli/v2"
)

func create(cliCtx *cli.Context) error {
	return retry(cliCtx, func(ctx *cli.Context) error {
		return bucketInstance.PutBucket(acl)
	})
}
