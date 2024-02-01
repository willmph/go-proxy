package main

import (
	"flag"
	"fmt"
	"io"
	"net"
)

func handleConnection(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	targetServer, err := net.Dial("tcp", targetAddr)
	if err != nil {
		fmt.Println("无法连接到目标服务器:", err)
		return
	}
	defer targetServer.Close()

	fmt.Printf("转发连接从 %s 到 %s\n", clientConn.RemoteAddr(), targetAddr)

	// 同时进行数据复制
	go func() {
		copyData(clientConn, targetServer)

	}()

	copyData(targetServer, clientConn)

}

func copyData(dst io.Writer, src io.Reader) {

	// 为了演示，这里简单地将数据读取到一个缓冲区，然后打印并写入目标
	buffer := make([]byte, 1024)
	for {
		n, err := src.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		// 打印数据
		fmt.Printf(" %s\n", string(buffer[:n]))

		// 将数据写入目标
		_, err = dst.Write(buffer[:n])
		if err != nil {
			return
		}
	}

}

func main() {
	listenPort := "55555" // 本地监听端口
	//targetServerAddr := "127.0.0.1:9090" // 目标服务器地址和端口

	targetServerAddr := flag.String("bind", "127.0.0.1:6000", "The address to bind to")
	flag.Parse()
	if *targetServerAddr == "" {
		fmt.Println("use default bind")
	}
	listen, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		fmt.Println("监听失败:", err)
		return
	}
	defer listen.Close()

	fmt.Printf("端口转发服务正在监听端口 %s，将转发到 %s\n", listenPort, *targetServerAddr)

	for {
		clientConn, err := listen.Accept()
		if err != nil {
			fmt.Println("接受连接失败:", err)
			continue
		}

		// 启动一个新的 goroutine 处理连接
		go handleConnection(clientConn, *targetServerAddr)
	}
}
