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
	"sort"
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
	acl sinastoragegosdk.ACL = sinastoragegosdk.Private
	// bucket实例
	bucketInstance *sinastoragegosdk.Bucket
	// 协议
	scheme string = "https://"
)

// 有效的acl列表
var aclMap = map[sinastoragegosdk.ACL]bool{
	sinastoragegosdk.Private:           true,
	sinastoragegosdk.PublicRead:        true,
	sinastoragegosdk.PublicReadWrite:   true,
	sinastoragegosdk.AuthenticatedRead: true,
}

// validateAcl 验证acl
func validateAcl(options map[string]string) error {
	acl = sinastoragegosdk.ACL(options["cannedAcl"])
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

var rootCmd = &cobra.Command{
	Use:   "SCS cli Tool",
	Short: "Cli Tool For SinaCloudStorage",
	Long:  `Cli Tool For SinaCloudStorage-SDK Build With Golang`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		accessKey = os.Getenv("S3_ACCESS_KEY_ID")
		if accessKey == "" {
			return errors.New("must set environment S3_ACCESS_KEY_ID")
		}

		secretKey = os.Getenv("S3_SECRET_ACCESS_KEY")
		if secretKey == "" {
			return errors.New("must set environment S3_SECRET_ACCESS_KEY")
		}

		hostnameEnv := os.Getenv("S3_HOSTNAME")
		if hostnameEnv != "" {
			hostname = hostnameEnv
		}

		if unencrypted {
			scheme = "http://"
		}

		s3 = sinastoragegosdk.NewSCS(accessKey, secretKey, scheme+hostname)
		return nil
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

// splitArgs 自定义解析变量
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
