package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func put(cliCtx *cli.Context) error {
	obj := cliCtx.String("object")
	filename := cliCtx.String("filename")
	return retry(cliCtx, func(ctx *cli.Context) error {
		multi, err := bucketInstance.InitMulti(obj)
		if err != nil {
			return err
		}
		// Todo: 默认先按照1GB一片分片上传，后续优化直接上传及分片计算
		partInfo, err := multi.PutPart(filename, acl, 1024*1024*1024)
		if err != nil {
			return err
		}
		listPart, err := multi.ListPart()
		if err != nil {
			return err
		}
		for k, v := range listPart {
			if partInfo[k].ETag != v.ETag {
				return fmt.Errorf("part not match")
			}
		}
		err = multi.Complete(listPart)
		if err != nil {
			return err
		}
		meta := make(map[string]string, 1)
		if expires := cliCtx.String("expires"); expires != "" {
			meta["x-sina-expire"] = expires
		}
		if errPutMeta := bucketInstance.PutMeta(obj, meta); errPutMeta != nil {
			fmt.Println(errPutMeta.Error())
		}
		return nil
	})
}
