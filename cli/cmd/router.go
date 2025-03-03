package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
	// 指定object
	object string
	// sdk实例
	s3 *sinastoragegosdk.SCS
	// 指定acl
	acl sinastoragegosdk.ACL
	// bucket实例
	bucketInstance *sinastoragegosdk.Bucket
	// 协议
	scheme = "https://"
	// 输出路径
	output *os.File
	// 是否打开了写文件
	hasWriter bool
)

// 命令定义
var app = &cli.App{
	Name:     "Cli Tool For Sina Cloud Storage Build With Golang",
	Usage:    "新浪云存储命令行工具-使用Golang构建",
	Version:  "1.0.0",
	Commands: router(),
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "access_key",
			Usage:    "the access key of Sina Cloud Storage",
			Aliases:  []string{"a"},
			Required: true,
			EnvVars:  []string{"S3_ACCESS_KEY_ID"},
		},
		&cli.StringFlag{
			Name:     "secret_key",
			Usage:    "the secret_key of Sina Cloud Storage",
			Aliases:  []string{"s"},
			Required: true,
			EnvVars:  []string{"S3_SECRET_ACCESS_KEY"},
		},
		&cli.StringFlag{
			Name:    "hostname",
			Usage:   "specify alternative host of Sina Cloud Storage",
			Aliases: []string{"n"},
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
		&cli.BoolFlag{
			Name:    "unencrypted",
			Usage:   "unencrypted (use HTTP instead of HTTPS)",
			Aliases: []string{"u"},
		},
		// &cli.BoolFlag{
		// 	Name:    "show-properties",
		// 	Usage:   "show response properties on stdout",
		// 	Aliases: []string{"s"},
		// },
		&cli.UintFlag{
			Name:        "retries",
			Usage:       "retry retryable failures this number of times (default is 5)",
			Aliases:     []string{"r"},
			Value:       5,
			Destination: &retries,
		},
	},
	Before: func(c *cli.Context) error {
		accessKey = c.String("access_key")
		if accessKey == "" {
			return errors.New("environment S3_ACCESS_KEY_ID or parameter access_key can't be empty")
		}

		secretKey = c.String("secret_key")
		if secretKey == "" {
			return errors.New("environment S3_SECRET_ACCESS_KEY or parameter secret_key can't be empty")
		}

		hostname = c.String("hostname")
		if hostname == "" {
			return errors.New("environment S3_HOSTNAME or parameter hostname can't be empty")
		}

		if c.Bool("unencrypted") {
			scheme = "http://"
		}

		s3 = sinastoragegosdk.NewSCS(accessKey, secretKey, scheme+hostname)
		fmt.Fprintf(os.Stdout, "[%s] Starting", time.Now().Format("20060102 15:04:05"))
		return nil
	},
	After: func(ctx *cli.Context) error {
		fmt.Fprintf(os.Stdout, "[%s] Finished", time.Now().Format("20060102 15:04:05"))
		return nil
	},
}

// 子命令列表
var config = map[string]*cli.Command{
	"create": {
		Usage:     "create a new bucket",
		ArgsUsage: "[bucket]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "cannedAcl",
				Usage: "Canned ACL for the bucket (see Canned ACLs)",
				Value: string(sinastoragegosdk.Private),
			},
		},
		Before: validateBucketObjectAcl,
		Action: create,
	},
	"delete": {
		Usage:     "delete a bucket or bucket/object",
		ArgsUsage: "[bucket[/object]]",
		Before:    validateBucketObject,
		Action:    delete,
	},
	"list": {
		Usage:     "lists owned buckets or list bucket contents",
		ArgsUsage: "[bucket]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "prefix",
				Usage: "prefix for results set",
			},
			&cli.StringFlag{
				Name:  "marker",
				Usage: "where in results set to start listing",
			},
			&cli.StringFlag{
				Name:  "delimiter",
				Usage: "delimiter for rolling up results set",
			},
			&cli.IntFlag{
				Name:  "maxkeys",
				Usage: "maximum number of keys to return in results set",
			},
		},
		Action: list,
	},
	"getacl": {
		Usage:     "get the ACL of a bucket or bucket/object",
		ArgsUsage: "[bucket[/object]]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "filename",
				Usage:       "output filename for ACL",
				DefaultText: "stdout",
			},
		},
		Before: validateBucketObjectWriter,
		Action: getacl,
	},
	"put": {
		Usage:     "puts an object",
		ArgsUsage: "[bucket/object]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "filename",
				Usage:       "filename to read source data from",
				Required:    true,
				DefaultText: "stdin",
			},
			&cli.StringFlag{
				Name:  "contentLength",
				Usage: "how many bytes of source data to put (required if source file is stdin)",
			},
			&cli.StringFlag{
				Name:  "cacheControl",
				Usage: "Cache-Control HTTP header string to associate with object",
			},
			&cli.StringFlag{
				Name:  "contentType",
				Usage: "Content-Type HTTP header string to associate with object",
			},
			&cli.StringFlag{
				Name:  "md5",
				Usage: "MD5 for validating source data",
			},
			&cli.StringFlag{
				Name:  "contentDispositionFilename",
				Usage: "Content-Disposition filename string to associate with object",
			},
			&cli.StringFlag{
				Name:  "contentEncoding",
				Usage: "Content-Encoding HTTP header string to associate with object",
			},
			&cli.StringFlag{
				Name:  "expires",
				Usage: "expiration date to associate with object",
			},
			&cli.StringFlag{
				Name:  "cannedAcl",
				Usage: "canned ACL for the object (see Canned ACLs)",
				Value: string(sinastoragegosdk.Private),
			},
			&cli.StringSliceFlag{
				Name:  "metadataHeaders",
				Usage: "metadata headers to associate with the object",
			},
		},
		Before: validateBucketObjectAcl,
		Action: put,
	},
	// Todo: 待实现
	"copy": {},
	"get": {
		Usage:     "Gets an object",
		ArgsUsage: "[bucket/object]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "filename",
				Usage:    "Filename to read source data from",
				Required: true,
			},
		},
		Before: validateBucketObject,
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
		Before: validateBucketObject,
		Action: head,
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

// validateBucketObject 验证bucket
func validateBucketObject(cliCtx *cli.Context) error {
	args := strings.SplitN(cliCtx.Args().First(), "/", 2)
	if len(args) == 1 {
		bucket = args[0]
	} else if len(args) == 2 {
		bucket = args[0]
		object = args[1]
	}
	if bucket == "" {
		return errors.New("parameter bucket can't be empty")
	}
	bucketInstance = s3.Bucket(bucket)
	return nil
}

// validateWriter 验证输出路径
func validateWriter(cliCtx *cli.Context) error {
	filename := cliCtx.String("filename")
	var err error
	if filename == "" {
		output = os.Stdout
	} else {
		if err = os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
			return err
		}
		output, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o755)
		if err != nil {
			return err
		}
		hasWriter = true
	}
	return nil
}

// validateBucketObjectAcl 验证bucket以及acl
func validateBucketObjectAcl(cliCtx *cli.Context) error {
	if err := validateBucketObject(cliCtx); err != nil {
		return err
	}
	if err := validateAcl(cliCtx); err != nil {
		return err
	}
	return nil
}

// validateBucketObjectWriter 验证bucket,bucket/object及filename
func validateBucketObjectWriter(cliCtx *cli.Context) error {
	if err := validateBucketObject(cliCtx); err != nil {
		return err
	}
	if err := validateWriter(cliCtx); err != nil {
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
	fmt.Fprintln(os.Stderr, err)
	if retries > 0 {
		var i uint
		for i = 1; i <= retries; i++ {
			fmt.Fprintf(os.Stdout, "Start retry %d\n", i)
			err := f(cliCtx)
			if err == nil {
				return nil
			}
		}
		fmt.Fprintln(os.Stderr, "All retries failed")
	}
	return nil
}
