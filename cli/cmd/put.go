/**
 * 上传1个对象
 * File  : cli/cmd/put.go
 * Author: JontyLee
 * Date  : 2025-02-25 12:03:24
 */
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	sinastoragegosdk "github.com/SinaCloudStorage/SinaCloudStorage-SDK-Go"
)

var putCmd = &cobra.Command{
	Use:     "put",
	Short:   "puts an object",
	Args:    cobra.ExactArgs(2),
	PreRunE: checkAcl,
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		bucketInstance = s3.Bucket(bucket)
		return retry(func() error {
			multi, err := bucketInstance.InitMulti(args[1])
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
			// meta := make(map[string]string, 1)
			// if expires := cliCtx.String("expires"); expires != "" {
			// 	meta["x-sina-expire"] = expires
			// }
			// if errPutMeta := bucketInstance.PutMeta(obj, meta); errPutMeta != nil {
			// 	fmt.Println(errPutMeta.Error())
			// }
			return nil
		})
	},
}

func init() {
	createCmd.Flags().StringVarP(&cannedAcl, "cannedAcl", "a", string(sinastoragegosdk.Private), "canned ACL for the bucket (see Canned ACLs)")
	rootCmd.AddCommand(putCmd)
}
