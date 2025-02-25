/**
 * cli处理
 * File  : cli/main.go
 * Author: JontyLee
 * Date  : 2025-02-17 14:34:10
 */
package main

import (
	"fmt"
	"os"
	"time"
)

// APP启动定义
func main() {
	fmt.Printf("[%s] Started\n", time.Now().Format("2006-01-02 15:04:05"))
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("[%s] Finished\n", time.Now().Format("2006-01-02 15:04:05"))
}
