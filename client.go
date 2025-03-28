package main

import (
	"net"
	"fmt"
	// "time"
	"flag"
	"io"
	"os"
)


type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int // 模式选择
}

// 创建一个client接口
func NewClient(serverIp string, serverPort int) *Client {
	// 创建一个client对象
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag: 999,
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

// 处理server回应的消息，直接显示到标准输出
func (client *Client) DealResponse() {
	// 一旦有消息，直接显示到标准输出
	// 监听服务器返回的消息
	io.Copy(os.Stdout, client.conn)
}

// 更新用户名业务
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>> 请输入用户名:")
	fmt.Scanln(&client.Name)
	sendMsg := fmt.Sprintf("rename|%s\n", client.Name)
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}

	return true
}

// 聊天模式菜单
func (client *Client) menu() bool {
	var flag int
	fmt.Println(">>>请选择聊天模式:")
	fmt.Println(">>>>> 1. 公聊模式")
	fmt.Println(">>>>> 2. 私聊模式")
	fmt.Println(">>>>> 3. 更新用户名")
	fmt.Println(">>>>> 0. 退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>> 请输入合法范围内的数字")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 { // 非退出
		for client.menu() != true { // 不合法输入
		}
		// 正常模式处理逻辑
		switch client.flag {
		case 1:
			fmt.Println(" 公聊模式选择...")
			break
		case 2:
			fmt.Println(" 私聊模式选择")
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器Ip地址（默认127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认8888）")
}

func main() {
	// 命令行解析
	flag.Parse()
	
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> 连接服务器失败")
		return
	}

	// 读取服务器返回消息
	go client.DealResponse()

	fmt.Println(">>>>> 连接服务器成功")

	// 启动客户端的业务
	client.Run()
}