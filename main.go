package main

import (
	"fmt"
	"net"
	"os"
)

func forward(conFrom, conTo net.Conn) {
	for {
		buff := make([]byte, 4096, 4096)
		n, err := conFrom.Read(buff)
		if err != nil {
			fmt.Println("con read fail", err)
			conFrom.Close()
			conTo.Close()
			return
		}
		_, err = conTo.Write(buff[:n])
		if err != nil {
			fmt.Println("con write fail", err)
			conFrom.Close()
			conTo.Close()
			return
		}
	}
}
func main() {
	arg_num := len(os.Args)
	if arg_num != 3 {
		fmt.Println("useage tcp-proxy [listen addr] [server addr]")
		return
	}
	listenaddr := os.Args[1]
	serveraddr := os.Args[2]
	ln, err := net.Listen("tcp", listenaddr)
	if err != nil {
		fmt.Printf("listen at [%s] fail", listenaddr, err)
	}
	fmt.Printf("tcp-proxy started  listen addr:[%s]  server addr:[%s] \n", listenaddr, serveraddr)
	for {
		client, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept fail", err)
			break
		}
		go func(client net.Conn) {
			proxy, err := net.Dial("tcp", serveraddr)
			if err == nil {
				fmt.Printf("new client:[%s] connected\n", client.RemoteAddr())
				go forward(client, proxy)
				go forward(proxy, client)
			} else {
				fmt.Println("connect to [%s] fail", serveraddr, err)
				client.Close()
			}
		}(client)
	}
}
