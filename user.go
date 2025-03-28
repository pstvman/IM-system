package main

import (
	"net"
	"strings"
)


type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn

	server *Server
}

// 创建一个User对象的接口
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,

		server: server,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 监听当前User channel的方法，一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <- this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

// 用户上线
func (this *User)Online() {
	//将用户加入到onlineMap中
	this.server.mapLock.Lock()
  this.server.OnlineMap[this.Name] = this
  this.server.mapLock.Unlock()

	// 广播当前用户上线消息
  this.server.BroadCast(this, "已上线")
}

// 用户下线
func (this *User)Offline() {
	this.server.mapLock.Lock()
  delete(this.server.OnlineMap, this.Name)
  this.server.mapLock.Unlock()

	// 广播当前用户下线消息
	this.server.BroadCast(this, "下线")
}

// 发送消息
func (this *User)SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 消息处理
func (this *User)DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			OnlineMsg := "[" + user.Addr + "]" + user.Name + ": 在线...\n"
			this.SendMsg(OnlineMsg)
		}
		this.server.mapLock.Unlock()
		
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 修改用户名
		// newName := strings.Split(msg, "|")[1]
		newName := msg[7:]

		// 判断name是否已经存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名已经被使用！\n")
			return
		}
		this.server.mapLock.Lock()
		delete(this.server.OnlineMap, this.Name)
		this.server.OnlineMap[newName] = this
		this.server.mapLock.Unlock()

		this.Name = newName
		this.SendMsg("您已经更新用户名：" + this.Name + "\n")

	} else if len(msg) >= 3 && msg[:3] == "to|" {
		// 私聊消息
		// 1. 获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("消息格式不正确，请使用\"to|用户名|消息内容\"格式\n")
			return
		}
		if remoteName == this.Name {
			this.SendMsg("不能给自己发送消息\n")
			return
		}

		// 2. 根据用户名得到对方User对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("该用户名不存在\n")
			return
		}

		// 3. 获取消息内容并发送给对方
		content := strings.Split(msg, "|")[2]
		remoteUser.SendMsg(this.Name + "对您说：" + content + "\n")

	} else {
		// 发送消息
		this.server.BroadCast(this, msg)
	}
}