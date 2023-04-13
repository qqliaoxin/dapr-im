package main

import (
	"bufio"
	"fmt"
	"net"
)

// 用户结构体
type Client struct {
	conn net.Conn
	name string
}

// 连接到客户端
func handleConn(conn net.Conn) {
	// 读取客户端消息
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg := input.Text()
		// 广播消息
		broadcast(msg, conn)
	}
}

// 广播消息
func broadcast(msg string, conn net.Conn) {
	// 获取所有在线用户
	clients := getClients()
	for client := range clients {
		if clients[client].conn != conn {
			// 除发送者外,向其他用户发送消息
			fmt.Fprintf(clients[client].conn, "%s: %s\n", clients[client].name, msg)
		}
	}
}

// 获取在线用户
func getClients() map[string]Client {
	clients := make(map[string]Client)
	// 遍历所有连接
	for conn, client := range clients {
		// 检查客户端是否断开
		if _, err := conn.(*net.TCPConn).Write([]byte{0}); err != nil {
			// 移除断开连接的客户端
			delete(clients, client.name)
			conn.Close()
		}
	}
	return clients
}

func main() {
	// 监听端口
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 运行服务器
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// 新连接加入
		go handleConn(conn)
	}
}
