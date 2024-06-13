package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"

	sinastoragegosdk "github.com/SinaCloudStorage/SinaCloudStorage-SDK-Go"
)

// 路由列表
var config = map[string]*cli.Command{
	"upload": {
		Usage: "上传文件",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "object",
				Usage:    "",
				Required: true,
				Aliases:  []string{"obj"},
			},
			&cli.StringFlag{
				Name:     "local_patch",
				Usage:    "",
				Required: true,
				Aliases:  []string{"lp"},
			},
			&cli.StringFlag{
				Name:      "acl",
				Usage:     "",
				Required:  true,
				Value:     string(sinastoragegosdk.Private),
				Validator: validateAcl,
			},
		},
		Action: func(ctx context.Context, cli *cli.Command) error {
			obj := cli.String("object")
			acl := sinastoragegosdk.ACL(cli.String("acl"))
			localPatch := cli.String("local_patch")
			err := newBucket(cli).Put(obj, localPatch, acl)
			if err != nil {
				fmt.Println(err)
			}
			return err
		},
	},
	"download": {
		Usage: "",
		Flags: []cli.Flag{},
		Action: func(context.Context, *cli.Command) error {
			return nil
		},
	},
}

// main 入口
func main() {
	// 定义入口
	cmd := &cli.Command{
		Name:  "SinaCloudStorage Command Line",
		Usage: "新浪云存储命令行",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "access_key",
				Usage:    "",
				Required: true,
				Aliases:  []string{"ak"},
			},
			&cli.StringFlag{
				Name:     "secret_key",
				Usage:    "",
				Required: true,
				Aliases:  []string{"sk"},
			},
			&cli.StringFlag{
				Name:     "end_point",
				Usage:    "",
				Required: true,
				Aliases:  []string{"ep"},
			},
			&cli.StringFlag{
				Name:     "bucket_name",
				Usage:    "",
				Required: true,
				Aliases:  []string{"bn"},
			},
		},
		Commands: router(),
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

// router 生成路由配置
func router() (result []*cli.Command) {
	for name, command := range config {
		command.Name = name
		result = append(result, command)
	}

	return result
}

// newBucket 创建新的Bucket
func newBucket(cli *cli.Command) *sinastoragegosdk.Bucket {
	ak := cli.String("access_key")
	sk := cli.String("secret_key")
	ep := cli.String("end_point")
	bn := cli.String("bucket_name")
	return sinastoragegosdk.NewSCS(ak, sk, ep).Bucket(bn)
}

// 有效的acl列表
var aclMap = map[sinastoragegosdk.ACL]bool{
	sinastoragegosdk.Private:           true,
	sinastoragegosdk.PublicRead:        true,
	sinastoragegosdk.PublicReadWrite:   true,
	sinastoragegosdk.AuthenticatedRead: true,
}

// validateAcl 验证acl
func validateAcl(acl string) error {
	if aclMap[sinastoragegosdk.ACL(acl)] {
		return nil
	}
	var aclList []string
	for aclValidate := range aclMap {
		aclList = append(aclList, string(aclValidate))
	}

	return errors.New(fmt.Sprintf("invalid acl: %s, acl must one of %s", acl, strings.Join(aclList, ",")))
}
