package main

import (
	"fmt"
	"time"
)

//进度条显示实现
type Bar struct {
	percent uint64
	cur 	uint64
	total	uint64
	rate 	string
	symbol  string
}

func (b *Bar)getPercent() {
	b.percent = uint64(float64(b.cur)/float64(b.total) * 100)
}

func (b *Bar)NewBarSym(start, total uint64, sym string) {
	b.cur = start
	b.total = total
	b.symbol = sym
	b.percent = uint64((float64(start)/float64(total)) * 100)
	for i:=0; uint64(i)<=b.percent; i+=2{
		b.rate += b.symbol
	}
}

func (b *Bar)NewBar(start, total uint64) {
	b.NewBarSym(start,total,"#")
}

func (b *Bar)Play(cur uint64) {
	oldPercent := b.percent
	b.cur = cur
	b.getPercent()
	if b.percent!=oldPercent && uint64(b.percent)%2==0{
		b.FreshRate()
		fmt.Printf("\r[%-50s]%3d%%   %8d/%d", b.rate, b.percent, b.cur, b.total)
	}
}

func (b* Bar)FreshRate()  {
	b.rate = ""
	for i:=1; uint64(i)<=b.percent; i+=2{
		b.rate += b.symbol
	}
}
func (b *Bar)Finish()  {
	fmt.Println()
}

func main() {
	var bar Bar
	bar.NewBarSym(0,100,"=")
	for i:=0; i<=100; i++{
		time.Sleep(time.Millisecond*200)
		bar.Play(uint64(i))
	}
	bar.Finish()
}

