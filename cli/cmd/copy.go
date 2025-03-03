package cmd

import (
	"github.com/urfave/cli/v2"
)

func copy(cliCtx *cli.Context) error {
	return retry(cliCtx, func(ctx *cli.Context) error {
		return bucketInstance.Copy(acl)
	})
}
