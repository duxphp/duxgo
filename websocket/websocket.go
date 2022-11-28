package websocket

import (
	"bytes"
	"encoding/json"
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/middleware"
	"github.com/gookit/event"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cast"
	"net/http"
	"sync"
	"time"
)

var Socket *Service

var (
	maxMsgSize = int64(512)
	pongWait   = 60 * time.Second
	newLine    = []byte{'\n'}
	space      = []byte{' '}
	pingPeriod = 10 * time.Second
)

// Init 默认初始化
func Init() {
	//defer pool.Release()
	Socket = New()
	Socket.Start()
}

func ReleaseSocket() {
	Socket.Pool.Release()
}

type Service struct {
	Socket     websocket.Upgrader
	Clients    map[string]Clients
	Users      map[string]Users
	Broadcast  chan *Broadcast
	Register   chan *Client
	Unregister chan *Client
	Pool       *ants.Pool
	Login      func(data string) (map[string]any, error)
}

type Client struct {
	Auth    string
	User    *User
	Socket  *websocket.Conn
	Mutex   sync.Mutex
	Send    chan *Message
	Message []byte
	Service *Service
}

type Users map[uint]*User
type Clients map[*Client]bool

type User struct {
	ID     uint
	Auth   string
	Client *Client
}

// Message is return msg
type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Broadcast struct {
	Client *Client
	Msg    []byte
}

func New() *Service {
	socket := websocket.Upgrader{}
	socket.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	pool, _ := ants.NewPool(200000)
	return &Service{
		Pool:       pool,
		Socket:     socket,
		Clients:    map[string]Clients{},
		Users:      map[string]Users{},
		Broadcast:  make(chan *Broadcast),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (r *Service) Handler(auth string) func(c echo.Context) error {
	return func(c echo.Context) error {
		var err error
		// 设置客户端信息
		var client Client
		client.Socket, err = r.Socket.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		//client.ID, _ = function.GetUuid()
		//if err != nil {
		//	return err
		//}
		client.Service = r
		client.Send = make(chan *Message)
		client.Auth = auth
		// 注册客户端
		r.Register <- &client

		r.Pool.Submit(func() {
			client.ServiceRead()
		})
		r.Pool.Submit(func() {
			client.ServiceWrite()
		})
		return nil
	}
}

// ServiceRead 获取客户端消息
func (c *Client) ServiceRead() {

	defer func() {
		c.Service.Unregister <- c
		c.Socket.Close()
	}()
	// SetReadLimit 设置对大致
	c.Socket.SetReadLimit(maxMsgSize)
	// SetReadDeadline 设置链接超时
	_ = c.Socket.SetReadDeadline(time.Now().Add(pongWait))

	c.Socket.SetPongHandler(func(appData string) error {
		//每次收到pong都把deadline往后推迟60秒
		c.Socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			// 错误处理
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				core.Logger.Debug().Err(err).Msg("websocket error")
			}
			break
		}

		message := bytes.TrimSpace(bytes.Replace(msg, newLine, space, -1))
		c.Service.Broadcast <- &Broadcast{
			Client: c,
			Msg:    message,
		}
	}
}

// ServiceWrite 写入客户端消息
func (c *Client) ServiceWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()
	for {
		select {
		// 写消息到当前的 websocket 连接
		case message, ok := <-c.Send:
			_ = c.Socket.SetWriteDeadline(time.Now().Add(pingPeriod))
			if !ok {
				// 关闭通道
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// NextWriter 为要发送的下一条消息返回一个写入器。写入器的Close方法将完整的消息刷新到网络。
			w, err := c.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			content, _ := json.Marshal(message)
			w.Write(content)
			// 将排队聊天消息添加到当前的 websocket 消息中
			n := len(c.Send)
			for i := 0; i < n; i++ {
				msg := <-c.Send
				content, _ = json.Marshal(msg)
				w.Write(newLine)
				w.Write(content)
			}
			if err := w.Close(); err != nil {
				return
			}
		//心跳保持
		case <-ticker.C:
			_ = c.Socket.SetWriteDeadline(time.Now().Add(pingPeriod))
			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Start 启动服务
func (r *Service) Start() {

	r.Pool.Submit(func() {
		for {
			select {
			case client := <-r.Register:
				// 注册 channel
				if r.Clients[client.Auth] == nil {
					r.Clients[client.Auth] = map[*Client]bool{}
				}
				r.Clients[client.Auth][client] = true
				client.SendMsg("coon", "successful connection to socket service")
			case client := <-r.Unregister:
				// 注销 channel
				client.Close()
			case data := <-r.Broadcast:
				// 广播 channel
				MessageStruct := Message{}
				err := json.Unmarshal(data.Msg, &MessageStruct)
				if err != nil {
					data.Client.SendMsg("err", "incorrect message format")
					continue
				}
				if MessageStruct.Type == "" {
					data.Client.SendMsg("err", "message type error")
					continue
				}
				switch MessageStruct.Type {
				// 登录
				case "login":
					var userInfo map[string]any
					// 鉴权数据获取
					if r.Login != nil {
						userInfo, err = r.Login(cast.ToString(MessageStruct.Data))
					} else {
						userInfo, err = middleware.NewJWT().ParsingToken(data.Client.Auth, cast.ToString(MessageStruct.Data))
					}
					if err != nil {
						data.Client.SendMsg("err", err.Error())
						data.Client.Close()
					}
					UID := cast.ToUint(userInfo["id"])
					if r.Users[data.Client.Auth] == nil {
						r.Users[data.Client.Auth] = map[uint]*User{}
					}
					user := &User{
						ID:     UID,
						Client: data.Client,
					}
					r.Users[data.Client.Auth][UID] = user
					data.Client.User = user
					data.Client.SendMsg("login", "登录成功")
				case "ping":
					data.Client.SendMsg("pong", "")
				default:
					if data.Client.User == nil {
						data.Client.SendMsg("error", "未授权登录")
						data.Client.Socket.WriteMessage(websocket.CloseMessage, []byte{})
						continue
					}
					r.Pool.Submit(func() {
						event.Fire("websocket."+MessageStruct.Type, map[string]any{
							"client":  data.Client,
							"message": &MessageStruct,
						})
					})
				}
			}
		}

	})
}

// Close 关闭通道
func (c *Client) Close() {
	close(c.Send)
	if _, ok := c.Service.Clients[c.Auth]; !ok {
		return
	}
	if _, ok := c.Service.Clients[c.Auth][c]; !ok {
		return
	}
	delete(c.Service.Clients[c.Auth], c)
	if _, ok := c.Service.Users[c.Auth]; !ok {
		return
	}
	if _, ok := c.Service.Users[c.Auth][c.User.ID]; !ok {
		return
	}
	delete(c.Service.Users[c.Auth], c.User.ID)
}

// SendMsg 发送消息
func (c *Client) SendMsg(Type string, message string, datas ...any) bool {
	var data any
	if len(datas) > 0 {
		data = datas[0]
	}
	select {
	case _, ok := <-c.Send:
		if !ok {
			return false
		}
	default:
		c.Send <- &Message{
			Type:    Type,
			Message: message,
			Data:    data,
		}
	}
	return true
}

// Event 事件对接
func Event(name string, call func(client *Client, message *Message) error) {
	event.On("websocket."+name, event.ListenerFunc(func(e event.Event) error {
		client := e.Get("client").(*Client)
		message := e.Get("message").(*Message)
		return call(client, message)
	}))
}
