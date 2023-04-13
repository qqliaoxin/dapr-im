package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// 客户端结构体
type Client struct {
	// websocket连接
	Conn *websocket.Conn
	// 发送消息的channel
	Send chan []byte
}

// 客户端连接
var clients = make(map[string]*Client)

// 消息广播channel
var broadcast = make(chan []byte)

// websocket消息处理
func wsMessage(conn *websocket.Conn) {
	// 获取客户端对象
	client := clients[conn.RemoteAddr().String()]
	// 监听消息发送channel
	go client.listenWrite()
	// 读取消息
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		// 广播消息
		broadcast <- msg
	}
}

// 监听消息发送channel
func (c *Client) listenWrite() {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	c.Conn.Close()
}

// 消息广播
func broadcastMessage() {
	for msg := range broadcast {
		// 广播给所有客户端
		for _, client := range clients {
			client.Send <- msg
		}
	}
}

func main() {
	// websocket路由
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			fmt.Println(err)
		}
		// 新客户端连接
		client := &Client{
			Conn: conn,
			Send: make(chan []byte),
		}
		clients[conn.RemoteAddr().String()] = client
		// 处理客户端消息
		go wsMessage(conn)
	})
	// 广播消息goroutine
	go broadcastMessage()
	fmt.Println("server start at :8000")
	http.ListenAndServe(":8000", nil)
}
