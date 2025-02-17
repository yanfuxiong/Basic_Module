package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"time"
)

func main() {

	//1、输出控制台打印进度条展示
	bar := progressbar.Default(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(50 * time.Millisecond)
		//log.Printf("Bytes:%d, Num:%d, Percent:%d", int64(bar.State().CurrentBytes), bar.State().CurrentNum, int8(bar.State().CurrentPercent))
	}

	//2、IO操作进度条控制台打印展示
	bar := progressbar.DefaultBytes(100, "下载中...")
	//bar := progressbar.DefaultBytes(-1, "处理中") // 未知长度的任务
	dstFp, _ := os.OpenFile("dst.gz", os.O_CREATE|os.O_WRONLY, 0644)
	srcFp, _ := os.OpenFile("go1.14.2.src.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	io.Copy(io.MultiWriter(dstFp, bar), srcFp)

	//3、 自定义进度条控制台打印展示
	bar := progressbar.NewOptions(1000,
		progressbar.OptionEnableColorCodes(true), // 启用颜色
		progressbar.OptionShowBytes(true),        // 显示字节单位
		progressbar.OptionSetWidth(15),           // 进度条宽度
		progressbar.OptionSetDescription("[cyan][1/3][reset] 处理中..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	for i := 0; i < 1000; i++ {
		bar.Add(1)
		time.Sleep(20 * time.Millisecond)
	}

	//4、IO操作使用应用自己的进度条更新逻辑
	bar := progressbar.DefaultBytes(100, "下载中...")
	dstFp, _ := os.OpenFile("dst.gz", os.O_CREATE|os.O_WRONLY, 0644)
	srcFp, _ := os.OpenFile("go1.14.2.src.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	// 使用 io.Copy 复制数据，并监控进度
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			
			GoUpdateProgressBar(bar.State().CurrentBytes, bar.State().Max) // 调用自定义进度更新函数
			if bar.State().CurrentBytes >= bar.State().Max {
				break
			}
		}
	}()
	_, err := io.Copy(io.MultiWriter(dstFp, bar), srcFp)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}
	io.Copy(io.MultiWriter(dstFp, bar), srcFp)
	bar.Finish()

}
