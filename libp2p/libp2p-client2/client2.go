package main

import (
	"bufio"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"io"
	"log"
	"os"
	"time"
)

const ProtocolID = "/libp2p/dcutr/text/1.0.0"

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

func handleStream(s network.Stream) {
	//rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go func() { //写
		for {
			writeStr := <-g_chWriteTxt
			//nSize, err := rw.Write([]byte(writeStr + "\n"))
			nSize, err := s.Write([]byte(writeStr + "\n"))
			if err != nil {
				fmt.Printf("rw.Write error :%+v", err)
				continue
			}
			fmt.Printf("Write string size %d success!", nSize)
		}
	}()

	go func() {
		buf := bufio.NewReader(s)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("Stream closed by peer")
					break
				}
				fmt.Printf("rw.ReadString error :%+v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			fmt.Printf("client2 received: %s", line)
		}
	}()
}

func main() {
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", 7889)) //固定端口

	node, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer node.Close()

	fmt.Println("Node2 Multiaddress:", node.Addrs())

	// 设置第一个节点的流处理器
	node.SetStreamHandler(ProtocolID, func(s network.Stream) {
		//rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		/*buf := bufio.NewReader(s)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			fmt.Printf("Node1 received: %s", line)
		}*/

		handleStream(s)
	})

	// 获取节点的 PeerInfo
	pid := node.ID()
	multAddr := node.Addrs()[0].String()

	fmt.Printf("\nSelf ID:%+v  ma:%+v \n", pid, multAddr)
	fmt.Printf("========================\n\n")

	DebugCmdLine()
	//select {}
}
