package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"

	sinastoragegosdk "github.com/SinaCloudStorage/SinaCloudStorage-SDK-Go"
)

var (
	// 新浪云存储access key
	accessKey string
	// 新浪云存储secret_key
	secretKey string
	// 指定新浪云存储域名
	hostname string
	// 重试次数
	retries uint
	// 指定bucket
	bucket string
	// sdk实例
	s3 *sinastoragegosdk.SCS
	// 指定acl
	acl sinastoragegosdk.ACL
	// bucket实例
	bucketInstance *sinastoragegosdk.Bucket
)

// 命令定义
var app = &cli.App{
	Name:     "Cli Tool For 新浪云存储 Build With Golang",
	Usage:    "新浪云存储命令行工具-使用Golang构建",
	Version:  "1.0.0",
	Commands: router(),
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "access_key",
			Usage:    "新浪云存储access key，优先读取此配置，如无则读取环境变量S3_ACCESS_KEY_ID",
			Aliases:  []string{"ak"},
			Required: true,
			EnvVars:  []string{"S3_ACCESS_KEY_ID"},
		},
		&cli.StringFlag{
			Name:     "secret_key",
			Usage:    "新浪云存储secret_key，优先读取此配置，如无则读取环境变量S3_SECRET_ACCESS_KEY",
			Aliases:  []string{"sk"},
			Required: true,
			EnvVars:  []string{"S3_SECRET_ACCESS_KEY"},
		},
		&cli.StringFlag{
			Name:    "hostname",
			Usage:   "指定新浪云存储域名，优先读取此配置，如无则读取环境变量S3_HOSTNAME",
			Aliases: []string{"hn"},
			EnvVars: []string{"S3_HOSTNAME"},
			Value:   "sinacloud.net",
		},
		// &cli.BoolFlag{
		// 	Name:    "force",
		// 	Usage:   "force operation despite warnings",
		// 	Aliases: []string{"f"},
		// },
		// &cli.BoolFlag{
		// 	Name:    "vhost-style",
		// 	Usage:   "use virtual-host-style URIs (default is path-style)",
		// 	Aliases: []string{"h"},
		// },
		// &cli.BoolFlag{
		// 	Name:    "unencrypted",
		// 	Usage:   "unencrypted (use HTTP instead of HTTPS)",
		// 	Aliases: []string{"u"},
		// },
		// &cli.BoolFlag{
		// 	Name:    "show-properties",
		// 	Usage:   "show response properties on stdout",
		// 	Aliases: []string{"s"},
		// },
		&cli.UintFlag{
			Name:    "retries",
			Usage:   "retry retryable failures this number of times (default is 5)",
			Aliases: []string{"r"},
			Value:   5,
		},
	},
	Before: func(c *cli.Context) error {
		accessKey = c.String("access_key")
		secretKey = c.String("secret_key")
		if accessKey == "" {
			return errors.New("must set environment S3_ACCESS_KEY_ID or parameter access_key")
		}
		if secretKey == "" {
			return errors.New("must set environment S3_SECRET_ACCESS_KEY or parameter secret_key")
		}
		hostname = c.String("hostname")
		if hostname == "" {
			return errors.New("hostname is empty")
		}
		retries = c.Uint("retries")
		s3 = sinastoragegosdk.NewSCS(accessKey, secretKey, hostname)
		return nil
	},
}

// 路由列表
var config = map[string]*cli.Command{
	"create": {
		Usage: "Create a new bucket",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bucket",
				Usage:    "Bucket to create",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "cannedAcl",
				Usage: "Canned ACL for the bucket (see Canned ACLs)",
				Value: string(sinastoragegosdk.Private),
			},
		},
		Before: validateBucketAcl,
		Action: create,
	},
	"delete": {
		Usage: "Delete a bucket or bucket/object",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bucket",
				Usage:    "Bucket to delete",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "object",
				Usage: "Bucket object to delete, must set parameter bucket at same time",
			},
		},
		Before: validateBucket,
		Action: delete,
	},
	"list": {
		Usage: "Lists owned buckets or list bucket contents",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "bucket",
				Usage: "Bucket to list",
			},
			&cli.StringFlag{
				Name:  "prefix",
				Usage: "Prefix for results set",
			},
			&cli.StringFlag{
				Name:  "marker",
				Usage: "Where in results set to start listing",
			},
			&cli.StringFlag{
				Name:  "delimiter",
				Usage: "Delimiter for rolling up results set",
			},
			&cli.StringFlag{
				Name:  "maxkeys",
				Usage: "Maximum number of keys to return in results set",
			},
		},
		Action: list,
	},
	"getacl": {
		Usage: "Get the ACL of a bucket or bucket/object",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bucket",
				Usage:    "Bucket to get the ACL of",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "object",
				Usage: "Bucket object to get the ACL of",
			},
			&cli.StringFlag{
				Name:  "filename",
				Usage: "Output filename for ACL (default is stdout)",
			},
		},
		Before: validateBucket,
		Action: getacl,
	},
	"put": {
		Usage: "Puts an object",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bucket",
				Usage:    "Bucket to put to",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "object",
				Usage:    "Bucket object to put to",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "filename",
				Usage:    "Filename to read source data from",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "cannedAcl",
				Usage: "Canned ACL for the bucket (see Canned ACLs)",
				Value: string(sinastoragegosdk.Private),
			},
			// Todo: 待实现
			// &cli.StringFlag{
			// 	Name:  "cacheControl",
			// 	Usage: "Cache-Control HTTP header string to associate with object",
			// },
			// &cli.StringFlag{
			// 	Name:  "contentType",
			// 	Usage: "Content-Type HTTP header string to associate with object",
			// },
		},
		Before: validateBucketAcl,
		Action: put,
	},
	// Todo: 待实现
	"copy": {},
	"get": {
		Usage: "Gets an object",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bucket",
				Usage:    "Bucket to put to",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "object",
				Usage:    "Bucket object to put to",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "filename",
				Usage:    "Filename to read source data from",
				Required: true,
			},
		},
		Before: validateBucket,
		Action: put,
	},
	"head": {
		Usage: "Gets only the headers of a bucket or bucket object",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bucket",
				Usage:    "Bucket to get the headers of",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "object",
				Usage: "Bucket object to get the headers of",
			},
		},
		Before: validateBucket,
	},
	// Todo: 待实现
	"gqs": {},
}

// router 生成路由配置
func router() (result []*cli.Command) {
	for name, command := range config {
		command.Name = name
		result = append(result, command)
	}
	return result
}

// 有效的acl列表
var aclMap = map[sinastoragegosdk.ACL]bool{
	sinastoragegosdk.Private:           true,
	sinastoragegosdk.PublicRead:        true,
	sinastoragegosdk.PublicReadWrite:   true,
	sinastoragegosdk.AuthenticatedRead: true,
}

// validateAcl 验证acl
func validateAcl(cliCtx *cli.Context) error {
	acl = sinastoragegosdk.ACL(cliCtx.String("cannedAcl"))
	if aclMap[acl] {
		return nil
	}
	var aclList []string
	for aclValidate := range aclMap {
		aclList = append(aclList, string(aclValidate))
	}
	sort.SliceStable(aclList, func(i, j int) bool {
		return aclList[i] < aclList[j]
	})
	return fmt.Errorf("invalid acl: %s, acl must one of %s", acl, strings.Join(aclList, ","))
}

// validateBucket 验证bucket
func validateBucket(cliCtx *cli.Context) error {
	bucket = cliCtx.String("bucket")
	if bucket == "" {
		return errors.New("must set parameter bucket")
	}
	bucketInstance = s3.Bucket(bucket)
	return nil
}

// validateBucketAcl 验证bucket以及acl
func validateBucketAcl(cliCtx *cli.Context) error {
	if err := validateBucket(cliCtx); err != nil {
		return err
	}
	if err := validateAcl(cliCtx); err != nil {
		return err
	}
	return nil
}

// retry 重试方法
func retry(cliCtx *cli.Context, f cli.ActionFunc) error {
	err := f(cliCtx)
	if err == nil {
		return nil
	}
	fmt.Println(err.Error())
	if retries > 0 {
		var i uint
		for i = 1; i <= retries; i++ {
			fmt.Printf("Start retry %d\n", i)
			err := f(cliCtx)
			if err == nil {
				return nil
			}
		}
		fmt.Println("All retries failed")
	}
	return nil
}
