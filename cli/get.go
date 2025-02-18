/**
 * Gets an object
 * File  : cli/get.go
 * Author: JontyLee
 * Date  : 2025-02-18 14:23:54
 */
package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/urfave/cli/v2"

	sinastoragegosdk "github.com/SinaCloudStorage/SinaCloudStorage-SDK-Go"
)

func get(cliCtx *cli.Context) error {
	channel := make(chan []string, 100)
	channelErr := make(chan *sinastoragegosdk.DownloadToChannelErr, 100)
	obj := cliCtx.String("object")
	filename := cliCtx.String("filename")
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	w := &sync.WaitGroup{}
	w.Add(1)
	go func() {
		defer w.Done()
		errRead := bucketInstance.DownloadToChannelWithErr(cliCtx.Context, obj, channel, channelErr, 100, '\n')
		if errRead != nil {
			fmt.Printf("%s\n", errRead.Error())
		}
	}()
	w.Add(1)
	go func() {
		defer w.Done()
		for errInfo := range channelErr {
			fmt.Printf("line_content:%s,error_msg:%s", errInfo.Err.Error(), errInfo.LineContent)
		}
	}()
	w.Add(1)
	go func() {
		defer w.Done()
		for data := range channel {
			for item := range data {
				_, errWrite := fmt.Fprintln(out, item)
				if errWrite != nil {
					fmt.Println(errWrite.Error())
				}
			}
		}
	}()
	w.Wait()
	return nil
}
