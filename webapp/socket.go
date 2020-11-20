package webapp

import (
	"metafora-bot-service/controller"
	"net/http"

	"golang.org/x/net/websocket"
	// "google.golang.org/api/drive/v3"
)

// ServerSoket Структура сокета
type ServerSoket struct {
	pattern    string
	clients    map[int]*Client
	addCh      chan *Client
	delCh      chan *Client
	sendAllCh  chan *Message
	doneCh     chan bool
	errCh      chan error
	Controller *controller.Controller
	// srv        *drive.Service
}

// NewServerSoket Создание структуры сокета
func NewServerSoket(pattern string, control *controller.Controller/*, srv *drive.Service*/) *ServerSoket {
	server := &ServerSoket{
		pattern,
		make(map[int]*Client),
		make(chan *Client),
		make(chan *Client),
		make(chan *Message),
		make(chan bool),
		make(chan error),
		control,
		// srv,
	}
	server.Controller.SetDeleteRoomClient(server.DeleteRoomClient)
	return server
}

// Add Добавление клиента
func (s *ServerSoket) Add(c *Client) {
	s.addCh <- c
}

// Del Удаление клиента
func (s *ServerSoket) Del(c *Client) {
	s.delCh <- c
}

// Send Отправка сообщения в канал сокета
func (s *ServerSoket) Send(msg *Message) {
	s.sendAllCh <- msg
}

// Done Завершение работы сокета
func (s *ServerSoket) Done() {
	s.doneCh <- true
}

// Err Ошибка сервера сокета
func (s *ServerSoket) Err(err error) {
	s.errCh <- err
}

// func (s *ServerSoket) sendPastMessages(c *Client) {
// 	for _, msg := range s.messages {
// 		c.Write(msg)
// 	}
// }

func (s *ServerSoket) sendAll(msg *Message) {
	for _, c := range s.clients {
		if len(c.auth) == 0 {
			Error("Клиент не авторизирован %v", msg)
			msg.Action = "Forbiden"
			msg.Body = "Клиент не авторизирован. Доступ запрещен"
			c.Write(msg)
			continue
		}
		// if msg.Action == "MessageInfo" {
		// 	if msg.ChatID == 10 {
		// 		c.Write(msg)
		// 	}
		// 	continue
		// }
		c.Write(msg)
	}
}

// WriteChats Отправка хука для клиента, что бы проверил соощения всем чатам
func (s *ServerSoket) WriteChats(auth string, msg *Message) {
	for _, c := range s.clients {
		if c.auth == auth {
			c.Write(msg)
		}
	}
}

// CheckClient Отправка хука для клиента, что бы проверил соощения
func (s *ServerSoket) CheckClient(auth string, chatID int64, message string, params interface{}) {
	for _, c := range s.clients {
		if c.auth == auth {
			msg := &Message{
				Action: "NewClientMessage",
				Body:   message,
				ChatID: chatID,
				Params: params,
			}
			c.Write(msg)
			return
		}
	}
}

// DeleteRoomClient Отправка хука для клиента, что бы удалил индификатор комнаты
func (s *ServerSoket) DeleteRoomClient(auth string, chatID int64, message string, params interface{}) {
	for _, c := range s.clients {
		if c.auth == auth {
			msg := &Message{
				Action: "DeleteRoom",
				Body:   message,
				ChatID: chatID,
				Params: params,
			}
			c.Write(msg)
			return
		}
	}
}

// ListenSocket Слушатель канала сокета
func (s *ServerSoket) ListenSocket() {
	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		client := NewClient(ws, s)
		s.Add(client)
		defer func() {
			s.Controller.ExitClientRoom(client.auth)
			s.Del(client)
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()
		s.Controller.SetOnMessage(s.CheckClient)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	Notice("Сервер сокета запущен по пути %s", s.pattern)

	for {
		select {
		// Add new a client
		case c := <-s.addCh:
			Info("Добавлен новый клиент")
			s.clients[c.id] = c
			Notice("Кол-во подключеных клиентов %d", len(s.clients))
			// del a client
		case c := <-s.delCh:
			Warning("Удаление клиента %d", c.id)
			delete(s.clients, c.id)
			// broadcast message for all clients
		case msg := <-s.sendAllCh:
			Info("Отправленно всем клиентам %v", msg)
			s.sendAll(msg)

		case err := <-s.errCh:
			Error("Ошибка сокета %s", err.Error())
			s.Done()

		case <-s.doneCh:
			return
		}
	}
}

// func (d *Resource) checkClientToken(client *Client) {
// 	if len(client.token) > 0 && client.auth == (AuthToken{}){
// 		auth := AuthToken{}
// 		token, _ := jwt.Parse(strings.Replace(client.token, "Bearer ", "", 1), func(token *jwt.Token) (interface{}, error) {

// 			// Don't forget to validate the alg is what you expect:
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 			}

// 			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
// 			return Token, nil
// 		})
// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			auth.IDAgents = ParseInt(claims["IDAgents"])
// 			auth.IdTerminal = ParseInt(claims["IdTerminal"])
// 			auth.IDSysUser = ParseInt(claims["IDSysUser"])
// 			auth.IDTypeTerminal = ParseInt(claims["IDTypeTerminal"])
// 			//Auth.Passmd5 = fmt.Sprintf("%s", claims["Passmd5"])
// 			//Auth.Password = claims["Password"].(string)
// 			auth.Sign = fmt.Sprintf("%s", claims["Sign"])
// 			auth.Guid = fmt.Sprintf("%s", claims["Guid"])
// 			auth.Created = ParseInt(claims["Created"])
// 			if claims["FiscalAuth"] != nil {
// 				if fiscal,ok := claims["FiscalAuth"].(map[string]interface{});ok {
// 					auth.FiscalAuth.Id = ParseInt64(fiscal["Id"])
// 					auth.FiscalAuth.Login = fmt.Sprintf("%s", fiscal["Login"])
// 					auth.FiscalAuth.UserType = (fiscal["UserType"])
// 					auth.FiscalAuth.FullName = fmt.Sprintf("%s", fiscal["FullName"])
// 					//Auth.FiscalAuth.Password = fmt.Sprintf("%s", claims["Password"])
// 				}
// 			}
// 			client.auth = auth
// 		}
// 	}
// }
