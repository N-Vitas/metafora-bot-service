package webapp

import (
	"database/sql"
	"fmt"
	"io"
	"metafora-bot-service/controller"
	"metafora-bot-service/controller/questions"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// Client Структура клиента
type Client struct {
	id     int
	ws     *websocket.Conn
	server *ServerSoket
	auth   string
	ch     chan *Message
	doneCh chan bool
	timer  time.Time
	doSell bool
}

const channelBufSize = 1024

var maxID = int(0)

// NewClient Создание нового клиента
func NewClient(ws *websocket.Conn, server *ServerSoket) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxID++
	return &Client{
		maxID,
		ws,
		server,
		"",
		make(chan *Message, channelBufSize),
		make(chan bool),
		time.Now(),
		false,
	}
}

// Conn Получение канала сокета
func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		msg := fmt.Sprintf("client %d is disconnected.", c.id)
		err := fmt.Errorf(msg)
		c.server.Err(err)
	}
}

// Done Завершение работы сокета
func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// listenWrite request via chanel
func (c *Client) listenWrite() {
	Notice("Слушатель сообщений для клиента")
	for {
		select {
		// send message to the client
		case msg := <-c.ch:
			websocket.JSON.Send(c.ws, msg)
			// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// SendFirstReplic Стартовая реплика бота
func (c *Client) SendFirstReplic(chatRoom string) {
	for c.doSell {
		if time.Now().After(c.timer) {
			var msg Message
			if room, ok := c.server.Controller.GetRoom(chatRoom); ok {
				c.doSell = false
				msg.Action = "Welcom"
				msg.Body = room.Replic.Message
				if c.server.Controller.SaveReplic(room) {
					room, _ = c.server.Controller.GetRoom(chatRoom)
					msg.Params = room
					c.WriteChats(&msg)
				}
			}
		}
		time.Sleep(2500 * time.Millisecond)
	}
}

// NextReplicBot Стартовая реплика бота
func (c *Client) NextReplicBot(room controller.ChatRoom, answer string) {
	msg := Message{} // Сообщение для клиента
	msg.Action = "Welcom"
	msg.Body = room.Replic.Message
	// Если реплика первая, то пропускаем проверку
	if room.Replic.ID == 1 && c.doSell {
		if c.server.Controller.SaveReplic(room) {
			room, _ = c.server.Controller.GetRoom(room.Room.Room)
			msg.Params = room
			msg.Body = room.Replic.Message
			c.WriteChats(&msg)
			c.doSell = false
			return
		}
		msg.Params = room
		msg.Body = room.Replic.Message
		c.WriteChats(&msg)
		c.doSell = false
		return
	}
	// Если у реплики нет ответа то просто шлем следующую
	if len(room.Replic.DataType) == 0 {
		if c.SendNextReplic(room, msg) {
			return
		}
		Error("Реплика не сохранилась")
		return
	}
	types := strings.Split(room.Replic.DataType, ",")
	for _, val := range types {
		if val == answer {
			if c.SendNextReplic(room, msg) {
				return
			}
			Error("Реплика не сохранилась")
			return
		}
	}
	if c.server.Controller.MaxSort() > room.Replic.Sort {
		if c.server.Controller.SaveReplic(room) {
			room, _ = c.server.Controller.GetRoom(room.Room.Room)
			msg.Params = room
			msg.Body = room.Replic.Message
			c.WriteChats(&msg)
		}
	}
}

// SendNextReplic Сохранение и отправка следующей реплики
func (c *Client) SendNextReplic(room controller.ChatRoom, msg Message) bool {
	next, err := questions.Next(room.Replic.Sort, c.server.Controller.GetTableName("bot"), c.server.Controller.DB.GetDb())
	if err != nil {
		if err == sql.ErrNoRows {
			Notice("%s", "Реплики закончились")
			return true
		}
		Error("Ошибка следующей реплики %s", err.Error())
		return false
	}
	room.Room.ReplicID = next.ID
	room.Replic = next
	if c.server.Controller.SaveReplic(room) {
		room, _ = c.server.Controller.GetRoom(room.Room.Room)
		msg.Params = room
		msg.Body = room.Replic.Message
		c.WriteChats(&msg)
		c.doSell = false
		return true
	}
	Error("%s", "Реплика не сохранилась")
	return false
}

// listenRead read request via chanel
func (c *Client) listenRead() {
	Notice("Слушатель сообщений от клиента")
	for {
		select {
		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return
			// read data from websocket connection
		default:
			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				if msg.Action == "Auth" {
					c.timer = c.server.Controller.NextReplicTime()
					c.auth = msg.Author
					c.server.Controller.ReopenRoom(msg.Author)
					msg.Body = "Welcom Chat Bot"
					c.doSell = true
					if room, ok := c.server.Controller.GetRoom(msg.Author); ok {
						msg.Params = room
						if len(room.Messages) > 0 {
							c.doSell = false
						}
					} else {
						msg.Action = "DeleteRoom"
						c.WriteChats(&msg)
						continue
					}
					c.WriteChats(&msg)
					go c.SendFirstReplic(c.auth)
					continue
				}
				if c.auth != msg.Author && len(c.auth) > 0 {
					msg.Action = "RefreshRoom"
					c.WriteChats(&msg)
					c.auth = ""
					continue
				}
				if msg.Action == "CreateMessage" {
					if room, ok := c.server.Controller.GetRoom(msg.Author); ok {
						c.doSell = false
						if room, ok = c.server.Controller.CreateClientMessage(msg.Body, room); ok {
							msg.Params = room
							msg.Action = "NewClientMessage"
							c.WriteChats(&msg)
							c.NextReplicBot(room, msg.Body)
						} else {
							msg.Action = "RefreshRoom"
							c.WriteChats(&msg)
						}
					}
				}
			}
		}
	}
}

// WriteChats отправка сообщений всем подключенным чатам с одним токеном
func (c *Client) WriteChats(msg *Message) {
	c.server.WriteChats(c.auth, msg)
}
