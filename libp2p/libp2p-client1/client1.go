package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"log"
	"os"
	"time"
)

const (
	ProtocolID  = "/libp2p/dcutr/text/1.0.0"
	PeerID      = "12D3KooWFFgo6b34vWx3DBvipwzhcUzYDUNBy4nEgHv7siQ9o9vG"
	PeermtuAddr = "/ip4/10.6.196.29/tcp/7889/p2p/"
)

var g_chWriteTxt = make(chan string)

func DebugCmdLine() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter text to debug:")
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("You entered:", line)
		g_chWriteTxt <- line
	}
}

func main() {
	node, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}
	defer node.Close()

	fmt.Printf("Node1 Multiaddress:%+v \n", node.Addrs())

	// 设置节点的流处理器
	node.SetStreamHandler(ProtocolID, func(s network.Stream) {
		buf := bufio.NewReader(s)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			fmt.Printf("client111 received: %s", line)
		}
	})

	// 获取节点的 PeerInfo
	pid := node.ID()
	multAddr := node.Addrs()[0].String()

	fmt.Printf("client1 Self ID:%+v  ma:%+v\n", pid, multAddr)
	fmt.Printf("========================\n\n")

	PeerMultiaddr, err := ma.NewMultiaddr(PeermtuAddr + PeerID)
	if err != nil {
		fmt.Printf("NewMultiaddr err")
		log.Fatalf("Failed to parse multiaddr: %v", err)
	}

	peerAddr, err := peer.AddrInfoFromP2pAddr(PeerMultiaddr)
	if err != nil {
		fmt.Printf("AddrInfoFromP2pAddr err")
		log.Fatalf("Failed to AddrInfoFromP2pAddr multiaddr: %v", err)
	}

	fmt.Printf("AddrInfoFromP2pAddr ok\n")
	// 连接对端节点
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//err = node.Connect(ctx, peer.AddrInfo{ID: PeerID, Addrs: []ma.Multiaddr{PeerMultiaddr}})
	err = node.Connect(ctx, *peerAddr)
	if err != nil { //连对端
		fmt.Printf("Connect err")
		log.Fatal(err)
	}

	fmt.Printf("Connected to peer ok\n")

	s, err := node.NewStream(ctx, peerAddr.ID, ProtocolID)
	if err != nil {
		fmt.Printf("NewStream err")
		return //log.Fatal(err)
	}

	nSize, err := s.Write([]byte("Hello xiongyanfu from Node1!\n"))
	if err != nil {
		fmt.Printf("Write err")
		return //log.Fatal(err)
	}
	fmt.Printf("s.Write size:%+v", nSize)

	go func() {
		for {
			readFromCmd := <-g_chWriteTxt
			s.Write([]byte(readFromCmd + "\n"))
			if err != nil {
				fmt.Printf("Write err")
				continue
			}
		}
	}()

	go func() {
		for {
			buf := bufio.NewReader(s)
			line, err := buf.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("Stream closed by peer")
					break
				}
				fmt.Printf("buf.ReadString error %+v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Printf("client1 received: %s", line)

		}
	}()

	defer s.Close()

	DebugCmdLine()
}
