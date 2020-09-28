package controller

import (
	"metafora-bot-service/controller/groups"
	"metafora-bot-service/controller/messages"
	"metafora-bot-service/controller/questions"
	"metafora-bot-service/controller/rooms"
	"metafora-bot-service/controller/settings"
	"strconv"
	"strings"
	"time"
)

// ChatRoom Расширенная структура комнаты
type ChatRoom struct {
	rooms.Room
	Group    groups.Group       `json:"groups"`
	Replic   questions.Replic   `json:"questions"`
	Messages []messages.Message `json:"messages"`
}

// GetRoom Регистрация комнаты
func (c *Controller) GetRoom(room string) (ChatRoom, bool) {
	r, e := rooms.Get(room, c.GetTableName("rooms"), c.DB.GetDb())
	// Получение существующей комнаты
	s := ChatRoom{r, groups.Group{}, questions.Replic{}, []messages.Message{}}
	if e != nil {
		// Создание комнаты
		if e = rooms.Create(room, c.GetTableName("rooms"), c.DB.GetDb()); e != nil {
			Error("Ошибка создания комнаты %s", e.Error())
			return s, false
		}
	}
	// Если комната есть собираем реплики и сообщения
	s.Replic, e = questions.Get(s.Room.ReplicID, c.GetTableName("bot"), c.DB.GetDb())
	if e != nil {
		Error("Ошибка загрузки реплилки %v", e)
	}
	// Возможно в комнате уже обозначена группа менеджеров
	s.Group, e = groups.Get(s.Room.GroupID, c.GetTableName("groups"), c.DB.GetDb())
	if e != nil {
		Error("Ошибка загрузки группы %v", e)
	}
	// Получение всех сообщений
	s.Messages, e = messages.GetAllInID(s.Room.MessagesID, c.GetTableName("message"), c.DB.GetDb())
	if e != nil {
		Error("Ошибка загрузки сообщений %v", e)
	}
	return s, true
}

// NextReplicTime Время для активации первой реплики бота
func (c *Controller) NextReplicTime() time.Time {
	// Берем настройки в базе
	s, e := settings.Get(c.GetTableName("settings"), c.DB.GetDb())
	if e == nil {
		Info("Время первой реплики бота в секундах %d", s.DurationStart)
		return time.Now().Add(time.Second * time.Duration(s.DurationStart))
	}
	Error("Ошибка настроек. По умолчанию первая реплика будет через 300 секунд %s", e.Error())
	return time.Now().Add(time.Second * time.Duration(300))
}

// SaveReplic Сохранение первой реплики бота
func (c *Controller) SaveReplic(room ChatRoom) bool {
	// Если есть менеджер
	if room.Room.ChatID > 0 {
		// Отправляем ему сообщение
		c.SendFromBot(room.Room.ChatID, room.Replic.Message, room.Room.ID)
	}
	if room.Room.ChatID == 0 && room.Room.GroupID > 0 {
		// Отправляем ему сообщение
		c.SendGroupFromBot(room.Room.GroupID, room.Replic.Message, room.Room.ID)
	}
	r, err := messages.NewMessage(c.GetTableName("message"), c.DB.GetDb(), "", room.Room.Room, room.Replic.Message, room.Replic.Type, room.Replic.DataType, -1, room.GroupID, room.ReplicID)
	if err != nil {
		Error("Ошибка создания сообщения %s", err.Error())
		return false
	}
	ids := []string{}
	if len(room.Room.MessagesID) > 0 {
		ids = strings.Split(room.Room.MessagesID, ",")
	}
	ids = append(ids, strconv.Itoa(int(r)))
	room.Room.LastMessage = r
	room.Room.MessagesID = strings.Join(ids, ",")
	err = rooms.Update(room.Room, c.GetTableName("rooms"), c.DB.GetDb())
	if err != nil {
		Error("Ошибка обновления комнаты %s", err.Error())
		return false
	}
	return true
}

// CreateClientMessage Сообщение от клиента
func (c *Controller) CreateClientMessage(message string, room ChatRoom) (ChatRoom, bool) {
	// Если есть менеджер
	if room.Room.ChatID > 0 {
		// Отправляем ему сообщение
		c.SendManager(room.Room.ChatID, message, room.Room.ID)
	}
	// Если нет группы
	if room.Room.GroupID == 0 {
		// Вытаскиваем все группы
		grs, err := groups.GetAll(c.GetTableName("groups"), c.DB.GetDb())
		if err != nil {
			Error("Ошибка получения всех групп %s", err.Error())
			return room, false
		}
		// перебираем группы
		for _, group := range grs {
			if strings.Index(message, group.Title) != -1 {
				room.Room.GroupID = group.ID
			}
		}
	}
	if room.Room.ChatID == 0 && room.Room.GroupID > 0 {
		// Отправляем ему сообщение
		c.SendGroup(room.Room.GroupID, message, room.Room.ID)
	}
	r, err := messages.NewMessage(c.GetTableName("message"), c.DB.GetDb(), "", room.Room.Room, message, "text", "", 0, room.Room.GroupID, room.Room.ReplicID)
	if err != nil {
		Error("Ошибка создания сообщения %s", err.Error())
		return room, false
	}
	ids := []string{}
	if len(room.Room.MessagesID) > 0 {
		ids = strings.Split(room.Room.MessagesID, ",")
	}
	ids = append(ids, strconv.Itoa(int(r)))
	room.Room.LastMessage = r
	room.Room.MessagesID = strings.Join(ids, ",")
	err = rooms.Update(room.Room, c.GetTableName("rooms"), c.DB.GetDb())
	if err != nil {
		Error("Ошибка обновления комнаты %s", err.Error())
		return room, false
	}
	// Получение всех сообщений
	room.Messages, err = messages.GetAllInID(room.Room.MessagesID, c.GetTableName("message"), c.DB.GetDb())
	if err != nil {
		Error("Ошибка загрузки сообщений %v", err)
	}
	return room, true
}
