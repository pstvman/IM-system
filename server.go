package main

import (
  "fmt"
  "net"
  "sync"
  "io"
  "time"
)


type Server struct {
  Ip    string
  Port  int

  //在线用户列表
  OnlineMap map[string]*User
  mapLock sync.RWMutex

  // 消息广播的channel
  Message chan string
}

// 创建一个Server的接口
func NewServer(ip string, port int) *Server {
  server := &Server{
    Ip: ip,
    Port: port,
    OnlineMap: make(map[string]*User),
    Message: make(chan string),
  }

  return server
}

// 启动服务器的接口
func (this *Server) Start() {
  // socket listen
  listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
  if err != nil {
    fmt.Println("net.Listen err: ", err)
    return
  }
  // close listen socket
  defer listener.Close()

  // 启动监听Message的goroutine
  go this.ListenMessager()

  for {
    // accept
    conn, err := listener.Accept()
    if err != nil {
      fmt.Println("listener accept err: ", err)
      continue
    }

    // do handler
    go this.Handler(conn)
  }
}

func (this *Server) Handler(conn net.Conn) {
  // 当前链接的业务
  // fmt.Println("链接建立成功")

  user := NewUser(conn, this)

  // 用户上线
  user.Online()

  //监听用户是否活跃的channel
  isLive := make(chan bool)

  // 接收客户端发送的消息
  go func() {
    buf := make([]byte, 4096)
    for {
      n, err := conn.Read(buf)
      if n == 0 {
        // 当前用户下线
        user.Offline()
        return
      }

      if err != nil && err != io.EOF {
        fmt.Println("Conn Read err: ", err)
        return
      }

      // 提取用户的消息，去除"\n"
      msg := string(buf[:n-1])

      // 消息处理
      user.DoMessage(msg)
      
      // 用户发消息，则说明在活跃
      isLive <- true
    }
  }()

  // 当前handler阻塞
  for {
    select {
    case <- isLive:
      // 当前用户是活跃的，应该重置定时器
      // 不做任何事情，为了激活select，更新下面的定时器

    case <- time.After(time.Second * 300): //定时器计时结束才向管道写数据
      // 超时踢登处理，将当前User强制关闭
      user.SendMsg("你被踢了\n")

      // 销毁用户管道资源
      close(user.C)

      // 关闭连接
      conn.Close()

      // 退出当前的handler
      return
    }
  }
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
  sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
  this.Message <- sendMsg
}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部在线User
func (this *Server) ListenMessager() {
  for {
    msg := <- this.Message

    // 将msg发送给全部在线User
    this.mapLock.Lock()
    for _, cli := range this.OnlineMap {
      cli.C <- msg
    }
    this.mapLock.Unlock()
  }
}