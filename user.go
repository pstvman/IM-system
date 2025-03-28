package main

import "net"


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

// 消息处理
func (this *User)DoMessage(msg string) {
		// 发送消息
		this.server.BroadCast(this, msg)
}