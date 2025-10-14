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

// context key（避免冲突）
type cancelSourceKey struct{}

// holder 保存本层的 source 以及指向父层 holder 的指针
type sourceHolder struct {
	mu     sync.Mutex
	src    CancelSource
	parent *sourceHolder
}

// WithCancelSource 返回一个带来源追踪的 ctx，以及一个接收 CancelSource 的 cancel 函数。
// 取消时会把来源写入 holder.src（首次写入有效），并调用底层 cancel 关闭 Done。
func WithCancelSource(parent context.Context) (context.Context, func(CancelSource)) {
	// 基于 parent 创建可取消 ctx
	base, cancel := context.WithCancel(parent)

	// 尝试从 parent 中拿到父 holder（若 parent 已经是由 WithCancelSource 创建）
	var parentHolder *sourceHolder
	if v := parent.Value(cancelSourceKey{}); v != nil {
		if ph, ok := v.(*sourceHolder); ok && ph != nil {
			parentHolder = ph
		}
	}

	// 创建本层 holder，并把 parentHolder 链接上
	h := &sourceHolder{
		parent: parentHolder,
	}

	// 把 holder 放到返回的 ctx 的 Value 中（注意：value 存的是指针，可被后续 cancel 更新）
	ctxWithVal := context.WithValue(base, cancelSourceKey{}, h)

	return ctxWithVal, func(src CancelSource) {
		// 首次写入 source
		h.mu.Lock()
		if h.src == 0 {
			h.src = src
		}
		h.mu.Unlock()
		// 触发 Done
		cancel()
	}
}

// GetCancelSource 从任意 ctx 中查找最近的 holder（通过 ctx.Value），并沿 parent 链向上查找第一个非零的 source。
// 返回 (source, true) 表示找到了来源，否则返回 (0, false)。
func GetCancelSource(ctx context.Context) (CancelSource, bool) {
	v := ctx.Value(cancelSourceKey{})
	if v == nil {
		return 0, false
	}
	h, ok := v.(*sourceHolder)
	if !ok || h == nil {
		return 0, false
	}
	// 沿 holder 链向上查找第一个非零 source
	for cur := h; cur != nil; cur = cur.parent {
		cur.mu.Lock()
		src := cur.src
		cur.mu.Unlock()
		if src != 0 {
			return src, true
		}
	}
	return 0, false
}

/* ---------- 演示 ---------- */
func worker(name string, ctx context.Context) {
	<-ctx.Done()
	if src, ok := GetCancelSource(ctx); ok {
		fmt.Printf("[%s] got cancel source: %v\n", name, src)
	} else {
		fmt.Printf("[%s] cancel but source unknown\n", name)
	}
}

func main() {
	// 根（祖父）ctx
	rootCtx, _ := WithCancelSource(context.Background())

	// 父 ctx（以 rootCtx 为 parent）
	parentCtx, parentCancel := WithCancelSource(rootCtx)

	// 子 ctx（以 parentCtx 为 parent）
	childCtx, _ := WithCancelSource(parentCtx)

	// 在子 goroutine 中等待 childCtx.Done()
	go worker("child", childCtx)
	go worker("parent", parentCtx)
	go worker("root", rootCtx)

	// 场景 1: 先由父级取消（来源 SourceUser）
	time.Sleep(100 * time.Millisecond)
	fmt.Println("== parent cancels with SourceUser ==")
	parentCancel(SourceUser)
	time.Sleep(200 * time.Millisecond)

	// 再复现：子层重新创建一套 ctx，验证祖父来源传递
	rootCtx2, rootCancel2 := WithCancelSource(context.Background())
	parentCtx2, _ := WithCancelSource(rootCtx2)
	childCtx2, _ := WithCancelSource(parentCtx2)
	go worker("child2", childCtx2)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("== root2 cancels with SourceSystem ==")
	rootCancel2(SourceSystem)
	time.Sleep(200 * time.Millisecond)

	// 最后示例：子自己被取消，优先级高于父或祖父
	rootCtx3, rootCancel3 := WithCancelSource(context.Background())
	parentCtx3, _ := WithCancelSource(rootCtx3)
	childCtx3, childCancel3 := WithCancelSource(parentCtx3)
	go worker("child3", childCtx3)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("== child3 cancels with SourceTimeout ==")
	childCancel3(SourceTimeout)
	// 即使祖父 later 取消，child3 已有来源
	time.Sleep(1 * time.Millisecond)
	rootCancel3(SourceSystem)
	time.Sleep(200 * time.Millisecond)
}
