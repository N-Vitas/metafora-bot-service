package restfull

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
)

// MessageService Формирование сервиса групп
func (app *Resource) MessageService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/message")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(app.GetAllMessage).
		Doc("Получение всех сообщений").
		Operation("GetAllGroup"))
	ws.Route(ws.GET("/{id}").To(app.GetMessageByIDRoom).
		Doc("Получение всех сообщений по ID комнаты").
		Operation("GetMessageByIDRoom").Param(ws.PathParameter("id", "ID").DataType("string")))
	ws.Route(ws.POST("/").To(app.GetMessageByID).
		Doc("Получение всех сообщений по ID").
		Operation("GetMessageByID"))
	return ws
}

// GetAllMessage Эндпоинт всех сообщений
func (app *Resource) GetAllMessage(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	messages := []Message{}
	query := "select id, img, chatRoom, message, chatID, groupID, replicID, status, date, type, dataType from " + app.Table("messages")
	rows, err := app.GetDb().Query(query)
	if err != nil {
		WriteStatusError(http.StatusForbidden, err, resp)
		return
	}
	for rows.Next() {
		ID := sql.NullInt64{}
		ChatID := sql.NullInt64{}
		GroupID := sql.NullInt64{}
		ReplicID := sql.NullInt64{}
		Status := sql.NullInt64{}
		Img := sql.NullString{}
		Room := sql.NullString{}
		Sqlmessage := sql.NullString{}
		Datetime := sql.NullString{}
		Type := sql.NullString{}
		DataType := sql.NullString{}
		err = rows.Scan(&ID, &Img, &Room, &Sqlmessage, &ChatID, &GroupID, &ReplicID, &Status, &Datetime, &Type, &DataType)
		if err != nil {
			Info("%v", err)
			continue
		}
		messages = append(messages, Message{
			ID:       ID.Int64,
			ChatID:   ChatID.Int64,
			GroupID:  GroupID.Int64,
			ReplicID: ReplicID.Int64,
			Status:   Status.Int64,
			Img:      Img.String,
			Room:     Room.String,
			Message:  Sqlmessage.String,
			Datetime: Datetime.String,
			Type:     Type.String,
			DataType: DataType.String,
		})
	}
	// Ответ пользователю
	WriteResponse(messages, resp)
}

// GetMessageByID Эндпоинт всех сообщений по ID комнаты
func (app *Resource) GetMessageByIDRoom(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	room := req.PathParameter("id")
	if len(room) == 0 {
		WriteStatusError(http.StatusBadRequest, errors.New("Нет паратеметра id комнаты"), resp)
		return

	}
	messages := []Message{}
	query := "select id, img, chatRoom, message, chatID, groupID, replicID, status, date, type, dataType from " + app.Table("messages") + " where chatRoom = ?"
	rows, err := app.GetDb().Query(query, room)
	if err != nil {
		WriteStatusError(http.StatusForbidden, err, resp)
		return
	}
	for rows.Next() {
		ID := sql.NullInt64{}
		ChatID := sql.NullInt64{}
		GroupID := sql.NullInt64{}
		ReplicID := sql.NullInt64{}
		Status := sql.NullInt64{}
		Img := sql.NullString{}
		Room := sql.NullString{}
		Sqlmessage := sql.NullString{}
		Datetime := sql.NullString{}
		Type := sql.NullString{}
		DataType := sql.NullString{}
		err = rows.Scan(&ID, &Img, &Room, &Sqlmessage, &ChatID, &GroupID, &ReplicID, &Status, &Datetime, &Type, &DataType)
		if err != nil {
			Info("%v", err)
			continue
		}
		messages = append(messages, Message{
			ID:       ID.Int64,
			ChatID:   ChatID.Int64,
			GroupID:  GroupID.Int64,
			ReplicID: ReplicID.Int64,
			Status:   Status.Int64,
			Img:      Img.String,
			Room:     Room.String,
			Message:  Sqlmessage.String,
			Datetime: Datetime.String,
			Type:     Type.String,
			DataType: DataType.String,
		})
	}
	// Ответ пользователю
	WriteResponse(messages, resp)
}

// GetMessageByID Эндпоинт всех сообщений по ID комнаты
func (app *Resource) GetMessageByID(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	msg := struct {
		MessagesID []int `json:"messagesID"`
	}{
		[]int{},
	}
	decoder := json.NewDecoder(req.Request.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		WriteStatusError(http.StatusBadRequest, errors.New("Не удалось распарсить данные"), resp)
		return
	}
	if len(msg.MessagesID) == 0 {
		WriteStatusError(http.StatusBadRequest, errors.New("Не верный список ID сообщений"), resp)
		return
	}
	list := []string{}
	for _, id := range msg.MessagesID {
		list = append(list, strconv.Itoa(id))
	}
	messages := []Message{}
	query := `select id, img, chatRoom, message, chatID, groupID, replicID, status, date, type, dataType from ` + app.Table("messages") + ` where id in(` + strings.Join(list, ",") + `)`
	rows, err := app.GetDb().Query(query)
	if err != nil {
		WriteStatusError(http.StatusForbidden, err, resp)
		return
	}
	for rows.Next() {
		ID := sql.NullInt64{}
		ChatID := sql.NullInt64{}
		GroupID := sql.NullInt64{}
		ReplicID := sql.NullInt64{}
		Status := sql.NullInt64{}
		Img := sql.NullString{}
		Room := sql.NullString{}
		Sqlmessage := sql.NullString{}
		Datetime := sql.NullString{}
		Type := sql.NullString{}
		DataType := sql.NullString{}
		err = rows.Scan(&ID, &Img, &Room, &Sqlmessage, &ChatID, &GroupID, &ReplicID, &Status, &Datetime, &Type, &DataType)
		if err != nil {
			Info("%v", err)
			continue
		}
		messages = append(messages, Message{
			ID:       ID.Int64,
			ChatID:   ChatID.Int64,
			GroupID:  GroupID.Int64,
			ReplicID: ReplicID.Int64,
			Status:   Status.Int64,
			Img:      Img.String,
			Room:     Room.String,
			Message:  Sqlmessage.String,
			Datetime: Datetime.String,
			Type:     Type.String,
			DataType: DataType.String,
		})
	}
	// Ответ пользователю
	WriteResponse(messages, resp)
}

// Message Структура Менеджера
type Message struct {
	ID       int64  `json:"id"`
	Img      string `json:"img"`
	Room     string `json:"chatRoom"`
	Message  string `json:"message"`
	ChatID   int64  `json:"chatID"`
	GroupID  int64  `json:"groupID"`
	ReplicID int64  `json:"replicID"`
	Status   int64  `json:"status"`
	Datetime string `json:"datetime"`
	Type     string `json:"type"`
	DataType string `json:"dataType"`
}
