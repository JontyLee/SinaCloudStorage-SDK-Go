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
)

// APP启动定义
func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
