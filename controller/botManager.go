package controller

import (
	"fmt"
	"metafora-bot-service/controller/groupmanager"
	"metafora-bot-service/controller/messages"
	"metafora-bot-service/controller/rooms"
	"metafora-bot-service/telegram"
	"strconv"
	"strings"
)

// SendManager Отправка сообщения менеджеру
func (c *Controller) SendManager(chatID int64, message string, roomID int64) {
	msg := fmt.Sprintf(`*Клиент написал*
	%s
	_Номер заявки %d_`, message, roomID)
	c.bot.Send(chatID, msg)
}

// SendFromBot Отправка сообщения менеджеру от бота
func (c *Controller) SendFromBot(chatID int64, message string, roomID int64) {
	msg := fmt.Sprintf(`*Робот написал клиенту*
	%s
	_Номер заявки %d_`, message, roomID)
	c.bot.Send(chatID, msg)
}

// SendGroup Отправка сообщения менеджеру
func (c *Controller) SendGroup(groupID int64, message string, roomID int64) {
	mngs, err := groupmanager.GetAll(groupID, c.GetTableName("managers"), c.GetTableName("group_manager"), c.DB.GetDb())
	if err != nil {
		Error("Ошибка поиска менеджеров по группе", err)
		return
	}
	for _, manager := range mngs {
		c.SendManager(manager.ChatID, message, roomID)
		Info("Отправлено менеджеру %d из группы %d в комнате %d", manager.ChatID, groupID, roomID)
	}
}

// SendGroupFromBot Отправка сообщения менеджеру от бота
func (c *Controller) SendGroupFromBot(groupID int64, message string, roomID int64) {
	mngs, err := groupmanager.GetAll(groupID, c.GetTableName("managers"), c.GetTableName("group_manager"), c.DB.GetDb())
	if err != nil {
		Error("Ошибка поиска менеджеров по группе", err)
		return
	}
	for _, manager := range mngs {
		c.SendFromBot(manager.ChatID, message, roomID)
		Info("Отправлено менеджеру %d из группы %d в комнате %d", manager.ChatID, groupID, roomID)
	}
}

// OnMessage Перехватчик сообщений от телеграмм-бота
func (c *Controller) OnMessage(chatID int64, message string) {
	if c.CheckComand(chatID, message) {
		return
	}
	room, err := rooms.FindByManager(chatID, c.GetTableName("rooms"), c.DB.GetDb())
	if err != nil {
		Warning("Сообщение %s. Комната с менеджером %d не обнаружена", message, chatID)
		c.bot.SendOpenButton(chatID, "Вы не выбрали номер чата")
		return
	}
	r, err := messages.NewMessage(c.GetTableName("message"), c.DB.GetDb(), "", room.Room, message, "text", "", chatID, room.GroupID, room.ReplicID)
	if err != nil {
		Error("Ошибка создания сообщения от менеджера %s", err.Error())
		return
	}
	ids := []string{}
	if len(room.MessagesID) > 0 {
		ids = strings.Split(room.MessagesID, ",")
	}
	ids = append(ids, strconv.Itoa(int(r)))
	room.LastMessage = r
	room.MessagesID = strings.Join(ids, ",")
	err = rooms.Update(room, c.GetTableName("rooms"), c.DB.GetDb())
	if err != nil {
		Error("Ошибка обновления комнаты %s", err.Error())
		return
	}
	params := ChatRoom{Room: room}
	// Получение всех сообщений
	params.Messages, err = messages.GetAllInID(room.MessagesID, c.GetTableName("message"), c.DB.GetDb())
	if err != nil {
		Error("Ошибка загрузки сообщений %v", err)
	}
	c.onMessage(room.Room, chatID, message, params)
}

// SetOnMessage Перехват уведомления клиента
func (c *Controller) SetOnMessage(onMessage func(chatRoom string, chatID int64, message string, params interface{})) {
	c.onMessage = onMessage
}

// CheckComand Проверка команд менеджера
func (c *Controller) CheckComand(chatID int64, message string) bool {
	switch message {
	case "Занять чат с клиентом":
		r, err := rooms.GetAllIsOpen(c.GetTableName("rooms"), c.DB.GetDb())
		if err != nil {
			c.bot.SendOpenButton(chatID, "Нет свободных чатов")
			Error("%v", err)
			return true
		}
		buttons := []int64{}
		for _, room := range r {
			if room.ChatID == 0 || room.ChatID == chatID {
				buttons = append(buttons, room.ID)
			}
		}
		if len(buttons) == 0 {
			c.bot.SendOpenButton(chatID, "Нет свободных чатов")
			return true
		}
		c.bot.RoomsButtonSend(chatID, buttons, message)
		return true
	case "Получить переписку чата":
		r, err := rooms.GetAllIsOpen(c.GetTableName("rooms"), c.DB.GetDb())
		if err != nil {
			c.bot.SendOpenButton(chatID, "Переписки по данному чату не найдено")
			Error("%v", err)
			return true
		}
		buttons := []int64{}
		for _, room := range r {
			buttons = append(buttons, room.ID)
		}
		if len(buttons) == 0 {
			c.bot.SendOpenButton(chatID, "Переписки по данному чату не найдено")
			return true
		}
		c.bot.RoomsButtonSend(chatID, buttons, message)
		return true
	case "Выйти из чата":
		err := rooms.ExitRoomMagager(chatID, c.GetTableName("rooms"), c.DB.GetDb())
		if err != nil {
			Error("CheckComand %v", err)
			c.bot.SendOpenButton(chatID, "У вас нет занятых чатов")
			return true
		}
		c.bot.SendOpenButton(chatID, "Вы освободили чат клиента")
		return true
	case "Закрыть чат":
		r, err := rooms.GetAllIsOpen(c.GetTableName("rooms"), c.DB.GetDb())
		if err != nil {
			c.bot.SendOpenButton(chatID, "У вас нет открытых чатов")
			Error("%v", err)
			return true
		}
		buttons := []int64{}
		for _, room := range r {
			if room.ChatID == chatID {
				buttons = append(buttons, room.ID)
			}
		}
		if len(buttons) == 0 {
			c.bot.SendOpenButton(chatID, "У вас нет открытых чатов")
			return true
		}
		c.bot.RoomsButtonSend(chatID, buttons, message)
		return true
	default:
		return false
	}
}

// OnComands Перехват команд от менеджера
func (c *Controller) OnComands(chatID int64, post telegram.PostRequests) {
	switch post.Message {
	case "Занять чат с клиентом":
		err := rooms.GoRoomMagager(post.ChatID, post.RoomID, c.GetTableName("rooms"), c.DB.GetDb())
		if err != nil {
			Error("OnComands %v", err)
			return
		}
		c.bot.SendOpenButton(chatID, fmt.Sprintf("Вы заняли чат номер %d", post.RoomID))
		return
	case "Получить переписку чата":
		msgs, err := rooms.GetMessagesRoom(post.RoomID, c.GetTableName("messages"), c.GetTableName("rooms"), c.DB.GetDb())
		if err != nil {
			Error("OnComands %v", err)
			return
		}
		filter := []string{}
		for _, message := range msgs {
			filter = append(filter, message.Message)
		}
		c.bot.SendOpenButton(chatID, fmt.Sprintf("*Переписка комнаты %d*\n\n%s", post.RoomID, strings.Join(filter, "\n")))
		return
	case "Закрыть чат":
		err := rooms.ClosedRoom(post.RoomID, c.GetTableName("rooms"), c.DB.GetDb())
		if err != nil {
			Error("OnComands %v", err)
			c.bot.SendOpenButton(chatID, "Закрытие чата отменено")
			return
		}
		c.bot.SendOpenButton(chatID, fmt.Sprintf("Вы закрыли чат номер %d", post.RoomID))
		// Отправка клиенту о закрытии комнаты
		c.DeleteRoomClient(post.RoomID)
		return
	default:
		return
	}
}

// DeleteRoomClient Функция для сокета чтоб клиент удалил индификатор комнаты
func (c *Controller) DeleteRoomClient(roomID int64) {
	room, err := rooms.GetByID(roomID, c.GetTableName("rooms"), c.DB.GetDb())
	if err != nil {
		Error("Комната %d не обнаружена", roomID)
		return
	}
	c.deleteRoomClient(room.Room, room.ChatID, "До свидания", room)
}

// SetDeleteRoomClient Функция для добавления функции сокета чтоб клиент удалил индификатор комнаты
func (c *Controller) SetDeleteRoomClient(f func(chatRoom string, chatID int64, message string, params interface{})) {
	c.deleteRoomClient = f
}
