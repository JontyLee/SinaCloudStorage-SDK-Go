/**
 * 主命令定义
 * File  : cmd/root.go
 * Author: JontyLee
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
	hostname string = "sinacloud.net"
	// 是否使用http代替https
	unencrypted bool
	// 重试次数
	retries uint
	// sdk实例
	s3 *sinastoragegosdk.SCS
	// 指定cannedAcl
	cannedAcl string
	// 使用的acl
	acl sinastoragegosdk.ACL
	// bucket实例
	bucketInstance *sinastoragegosdk.Bucket
	// 协议
	scheme string = "https://"
	// 查询前缀
	prefix string
	// 查询开始位置
	marker string
	// 分隔符设置
	delimiter string
	// 最大返回数量
	maxKeys int
	// 本地文件
	filename string
)

// 有效的acl列表
var aclMap = map[sinastoragegosdk.ACL]bool{
	sinastoragegosdk.Private:           true,
	sinastoragegosdk.PublicRead:        true,
	sinastoragegosdk.PublicReadWrite:   true,
	sinastoragegosdk.AuthenticatedRead: true,
}

// validateAcl 验证acl
func validateAcl(acl sinastoragegosdk.ACL) error {
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
	Use:   "s3",
	Short: "Cli Tool For SinaCloudStorage-SDK",
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
	rootCmd.Flags().UintVarP(&retries, "retries", "r", 5, "retry retryable failures this number of times")
	rootCmd.Flags().BoolVarP(&unencrypted, "unencrypted", "u", false, "use HTTP instead of HTTPS")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// retry 重试方法
func retry(f func() error) error {
	err := f()
	if err == nil {
		return nil
	}
	fmt.Fprintln(os.Stderr, err)
	var i uint
	for i = 0; i < retries; i++ {
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

// checkAcl 检查acl是否符合预期
func checkAcl(cmd *cobra.Command, args []string) error {
	acl = sinastoragegosdk.ACL(cannedAcl)
	if err := validateAcl(acl); err != nil {
		return err
	}
	return nil
}
