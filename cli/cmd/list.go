/**
 * List owned buckets or list bucket contents
 * File  : cli/list.go
 * Author: JontyLee
 * Date  : 2025-02-18 11:09:36
 */
package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func list(cliCtx *cli.Context) error {
	fmt.Fprintf(os.Stdout, "args: %+v\n", cliCtx.Args().Slice())
	bucket := cliCtx.Args().First()
	bucketInstance = s3.Bucket(bucket)
	if bucket == "" {
		return retry(cliCtx, func(ctx *cli.Context) error {
			data, err := bucketInstance.ListBucket()
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "%s\n", data)
			return nil
		})
	}
	prefix := cliCtx.String("prefix")
	delimiter := cliCtx.String("delimter")
	marker := cliCtx.String("marker")
	maxKeys := cliCtx.Int("maxkeys")
	return retry(cliCtx, func(ctx *cli.Context) error {
		data, err := bucketInstance.ListObject(prefix, delimiter, marker, maxKeys)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%s\n", data)
		return nil
	})
}
