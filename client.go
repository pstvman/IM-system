package main

import (
	"net"
	"fmt"
	"time"
)


type client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
}

// 创建一个client接口
func NewClient(serverIp string, serverPort int) *client {
	// 创建一个client对象
	client := &client{
		ServerIp: serverIp,
		ServerPort: serverPort,
	}

	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return nil
	}
	client.conn = conn

	// 返回对象
	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>> 连接服务器失败")
		return
	}

	fmt.Println(">>>>> 连接服务器成功")

	// 启动客户端的业务
	// select {}
	for {
		time.Sleep(time.Second)
	}
}