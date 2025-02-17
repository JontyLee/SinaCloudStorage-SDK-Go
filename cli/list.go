package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func list(cliCtx *cli.Context) error {
	bucket = cliCtx.String("bucket")
	bucketInstance = s3.Bucket(bucket)
	if bucket == "" {
		return retry(cliCtx, func(ctx *cli.Context) error {
			data, err := bucketInstance.ListBucket()
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", data)
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
		fmt.Printf("%s\n", data)
		return nil
	})
}
