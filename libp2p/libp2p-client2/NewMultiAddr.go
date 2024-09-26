package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	// 创建第一个节点
	h1, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}
	defer h1.Close()

	// 打印第一个节点的地址
	fmt.Println("Node1 Multiaddress:", h1.Addrs())

	// 创建第二个节点
	h2, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}
	defer h2.Close()

	// 打印第二个节点的地址
	fmt.Println("Node2 Multiaddress:", h2.Addrs())

	// 设置第一个节点的流处理器
	h1.SetStreamHandler("/text/1.0.0", func(s network.Stream) {
		buf := bufio.NewReader(s)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			fmt.Printf("Node1 received: %s", line)
		}
	})

	// 设置第二个节点的流处理器
	h2.SetStreamHandler("/text/1.0.0", func(s network.Stream) {
		buf := bufio.NewReader(s)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			fmt.Printf("Node2 received: %s", line)
		}
	})

	// 获取第一个节点的 PeerInfo
	pid1 := h1.ID()
	ma1 := h1.Addrs()[0]

	// 获取第二个节点的 PeerInfo
	pid2 := h2.ID()
	ma2 := h2.Addrs()[0]
	fmt.Printf("node1 id:%s ma:%+v,  node2 id:%s ma:%+v", pid1, ma1, pid2, ma2)
	// 连接两个节点
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h2.Connect(ctx, peer.AddrInfo{ID: pid1, Addrs: []ma.Multiaddr{ma1}}); err != nil {
		log.Fatal(err)
	}

	// 从第二个节点向第一个节点发送消息
	s, err := h2.NewStream(ctx, pid1, "/text/1.0.0")
	if err != nil {
		log.Fatal(err)
	}
	_, err = s.Write([]byte("Hello xiongyanfu data from Node2!\n"))
	if err != nil {
		log.Fatal(err)
	}
	s.Close()

	// 等待消息处理
	time.Sleep(1 * time.Second)
}
