package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	var target string
	flag.StringVar(&target, "target", "", "")
	flag.StringVar(&target, "t", "", "")

	var port string
	flag.StringVar(&port, "port", "443", "")
	flag.StringVar(&port, "p", "443", "")

	flag.Parse()
	if target == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -t target -p port \n", os.Args[0])
		return
	}

	// 建立 TCP 连接
	conn, err := net.Dial("tcp", target+":"+port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to server: %s\n", err)
		return
	}
	defer conn.Close()
	//建立TLS连接
	// 创建 TLS 连接
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName: target,
	})

	// 手动进行 TLS 握手
	if err := tlsConn.Handshake(); err != nil {
		fmt.Fprintf(os.Stderr, "TLS handshake error: %s\n", err)
		return
	}

	// 从标准输入中 ---> tlsConn
	go func() {
		_, err = io.Copy(tlsConn, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "io.Copy failed,err: %s\n", err)
			return
		}
	}()

	// tlsConn ---> 标准输出
	_, err = io.Copy(os.Stdout, tlsConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "io.Copy failed,err: %s\n", err)
		return
	}

}
