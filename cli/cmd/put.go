package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func put(cliCtx *cli.Context) error {
	filename := cliCtx.String("filename")
	return retry(cliCtx, func(ctx *cli.Context) error {
		multi, err := bucketInstance.InitMulti(object)
		if err != nil {
			return err
		}
		// Todo: 默认先按照每片1GB上传，后续优化根据文件大小区分直接上传还是分片上传
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
		if errPutMeta := bucketInstance.PutMeta(object, meta); errPutMeta != nil {
			fmt.Println(errPutMeta.Error())
		}
		return nil
	})
}
