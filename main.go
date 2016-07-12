package main

import (
	"fmt"
	"net"
	"os"
)

type proxy struct {
	id              int64
	client, backend net.Conn
}

func forward(conFrom, conTo net.Conn, uid int64) {
	for {
		buff := make([]byte, 4096, 4096)
		n, err := conFrom.Read(buff)
		if err != nil {
			fmt.Printf("[%d-%s] read fail:%s\n", uid, conFrom.RemoteAddr(), err.Error())
			conFrom.Close()
			conTo.Close()
			return
		}
		_, err = conTo.Write(buff[:n])
		if err != nil {
			fmt.Printf("[%d-%s] write fail:%s\n", uid, conFrom.RemoteAddr(), uid, err.Error())
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
		var uid int64
		uid = 1
		go func(client net.Conn, uid int64) {
			proxy, err := net.Dial("tcp", serveraddr)

			if err == nil {
				fmt.Printf("new client:[%d-%s] connected\n", uid, client.RemoteAddr())
				go forward(client, proxy, uid)
				go forward(proxy, client, uid)

			} else {
				fmt.Println("connect to [%s] fail", serveraddr, err)
				client.Close()
			}

		}(client, uid)
		uid++
	}
}
