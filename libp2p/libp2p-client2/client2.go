package main

import (
	"bufio"
	"encoding/pem"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"io"
	"log"
	"os"
	"time"
)

// const ProtocolID = "/libp2p/dcutr/direct/1.0.0"
const ProtocolID = "/echo/1.0.0"

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

func MarshalPrivateKeyToPEM(key crypto.PrivKey) ([]byte, error) {
	encoded, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %v", err)
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: encoded,
	})
	return pemEncoded, nil
}

func UnmarshalPrivateKeyFromPEM(pemData []byte) (crypto.PrivKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}
	return crypto.UnmarshalPrivateKey(block.Bytes)
}
func GenKey() crypto.PrivKey {
	privKeyFile := ".priv.pem"
	var priv crypto.PrivKey
	var err error
	var content []byte
	content, err = os.ReadFile(privKeyFile)
	if err != nil {
		priv, _, err = crypto.GenerateKeyPair(crypto.RSA, 2048)
		if err != nil {
			log.Fatal(err)
		}

		jsonData, err := MarshalPrivateKeyToPEM(priv)
		err = os.WriteFile(privKeyFile, jsonData, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return priv
	}

	priv, err = UnmarshalPrivateKeyFromPEM(content)

	return priv
}

func writeNodeID(ID string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(ID))
	if err != nil {
		log.Println(err)
	}
}
func main() {
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", 7889)) //固定端口

	priv := GenKey()

	node, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(priv),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer node.Close()

	writeNodeID(node.ID().String(), ".ID")

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
