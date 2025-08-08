package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 定义取消原因类型（标识不同调用方）
type CancelSource int

const (
	SourceUser CancelSource = iota + 1
	SourceSystem
	SourceTimeout
)

// 自定义Context（关键结构）
type customContext struct {
	context.Context
	mu *sync.Mutex
	//cancel context.CancelFunc
	source CancelSource // 记录取消来源
}

// 创建可识别来源的Context
func WithCancelSource(parent context.Context) (context.Context, func(CancelSource)) {
	ctx, cancel := context.WithCancel(parent)
	cc := &customContext{
		Context: ctx,
		mu:      &sync.Mutex{},
		//cancel:  nil,
		source: 0,
	}

	// 返回自定义取消函数（可携带来源参数）
	return cc, func(source CancelSource) {
		cc.mu.Lock()
		defer cc.mu.Unlock()
		if cc.source == 0 { // 只记录首次取消来源
			cc.source = source
		}
		cancel() // 触发标准取消
	}
}

// 获取取消来源（业务goroutine调用）
func GetCancelSource(ctx context.Context) (CancelSource, bool) {
	if cc, ok := ctx.(*customContext); ok {
		cc.mu.Lock()
		defer cc.mu.Unlock()
		return cc.source, cc.source != 0
	}
	return 0, false
}

// 示例业务逻辑
func businessWorker(ctx context.Context) {
	select {
	case <-ctx.Done():
		if source, ok := GetCancelSource(ctx); ok {
			switch source {
			case SourceUser:
				fmt.Println("用户手动取消")
			case SourceSystem:
				fmt.Println("系统维护取消")
			case SourceTimeout:
				fmt.Println("操作超时取消")
			}
		}
	case <-time.After(5 * time.Second):
		fmt.Println("业务正常完成")
	}
}

func main() {
	fmt.Printf("start")
	startTime := time.Now().UnixMilli()
	// 创建自定义Context
	ctx, cancel := WithCancelSource(context.Background())

	// 启动业务goroutine
	go businessWorker(ctx)

	// 模拟不同取消来源（示例）
	cancel(SourceSystem) // 传递取消来源

	time.Sleep(1 * time.Second)
	fmt.Printf("use[%d]ms", time.Now().UnixMilli()-startTime)
}
