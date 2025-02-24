/**
 * 主命令定义
 * File  : cmd/root.go
 * Author: jianlin6
 * Date  : 2025-02-24 14:55:00
 */
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	sinastoragegosdk "github.com/SinaCloudStorage/SinaCloudStorage-SDK-Go"
)

var (
	// 新浪云存储access key
	accessKey string
	// 新浪云存储secret_key
	secretKey string
	// 指定新浪云存储域名
	hostname string
	// 是否使用http代替https
	unencrypted bool
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
	// 协议
	scheme string = "https://"
)

var rootCmd = &cobra.Command{
	Use:   "SCS cli Tool",
	Short: "Cli Tool For SinaCloudStorage",
	Long:  `Cli Tool For SinaCloudStorage-SDK Build With Golang`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		accessKey = os.Getenv("S3_ACCESS_KEY_ID")
		if accessKey == "" {
			fmt.Fprintln(os.Stderr, errors.New("must set environment S3_ACCESS_KEY_ID"))
			os.Exit(1)
		}

		secretKey = os.Getenv("S3_SECRET_ACCESS_KEY")
		if secretKey == "" {
			fmt.Fprintln(os.Stderr, errors.New("must set environment S3_SECRET_ACCESS_KEY"))
			os.Exit(1)
		}

		hostname = os.Getenv("S3_HOSTNAME")
		if hostname == "" {
			fmt.Fprintln(os.Stderr, errors.New("must set environment S3_HOSTNAME"))
			os.Exit(1)
		}

		if unencrypted {
			scheme = "http://"
		}

		s3 = sinastoragegosdk.NewSCS(accessKey, secretKey, scheme+hostname)
	},
}

func init() {
	rootCmd.LocalFlags().UintVarP(&retries, "retries", "r", 5, "retry retryable failures this number of times")
	rootCmd.LocalFlags().BoolVarP(&unencrypted, "unencrypted", "u", false, "use HTTP instead of HTTPS")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// retry 重试方法
func retry(f func() error) error {
	var i uint
	var err error
	for i = 0; i <= retries; i++ {
		fmt.Printf("Start retry %d\n", i)
		err = f()
		if err == nil {
			return nil
		}
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("All retries failed")
	return err
}

func splitArgs(args []string, startIndex int) map[string]string {
	res := make(map[string]string, len(args))
	for i := startIndex; i < len(args); i++ {
		if args[i] == "" {
			continue
		}
		argSlice := strings.Split(args[i], "=")
		if len(argSlice) > 1 && argSlice[0] != "" && argSlice[1] != "" {
			res[argSlice[0]] = res[argSlice[1]]
		}
	}
	return res
}
