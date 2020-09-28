package restfull

import (
	"database/sql"
	"net/http"

	"github.com/emicklei/go-restful"
)

// RoomService Формирование сервиса групп
func (app *Resource) RoomService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/room")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(app.GetAllRoom).
		Doc("Получение всех комнат").
		Operation("GetAllRoom"))
	return ws
}

// GetAllRoom Эндпоинт всех сообщений
func (app *Resource) GetAllRoom(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	rooms := []Room{}
	query := `select r.id, r.chatID, r.groupID, r.replicID, r.lastmessage, r.messagesID, r.chatRoom, r.mute, r.status, r.date, g.title as groupTitle 
	from ` + app.Table("rooms") + ` r left join ` + app.Table("groups") + ` g on g.id = r.groupID`
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
		LastMessage := sql.NullInt64{}
		MessagesID := sql.NullString{}
		Sqlroom := sql.NullString{}
		Mute := sql.NullBool{}
		Status := sql.NullInt64{}
		Datetime := sql.NullString{}
		GroupTitle := sql.NullString{}

		err = rows.Scan(&ID, &ChatID, &GroupID, &ReplicID, &LastMessage, &MessagesID, &Sqlroom, &Mute, &Status, &Datetime, &GroupTitle)
		if err != nil {
			Info("%v", err)
			continue
		}
		name := func(title sql.NullString) string {
			if title.Valid {
				return title.String
			}
			return "Нет группы"
		}
		rooms = append(rooms, Room{
			ID:          ID.Int64,
			ChatID:      ChatID.Int64,
			GroupID:     GroupID.Int64,
			ReplicID:    ReplicID.Int64,
			LastMessage: LastMessage.Int64,
			MessagesID:  MessagesID.String,
			Room:        Sqlroom.String,
			Mute:        Mute.Bool,
			Status:      Status.Int64,
			Datetime:    Datetime.String,
			GroupTitle:  name(GroupTitle),
		})
	}
	// Ответ пользователю
	WriteResponse(rooms, resp)
}

// Room Структура Менеджера
type Room struct {
	ID          int64  `json:"id"`
	ChatID      int64  `json:"chatID"`
	GroupID     int64  `json:"groupID"`
	ReplicID    int64  `json:"replicID"`
	LastMessage int64  `json:"lastmessage"`
	MessagesID  string `json:"messagesID"`
	Room        string `json:"chatRoom"`
	Mute        bool   `json:"mute"`
	Status      int64  `json:"status"`
	Datetime    string `json:"date"`
	GroupTitle  string `json:"groupTitle"`
}
