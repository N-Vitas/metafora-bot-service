package crontimes

import (
	"database/sql"
	"fmt"
	"metafora-bot-service/controller/groups"
	"metafora-bot-service/controller/messages"
	"metafora-bot-service/controller/rooms"
	"time"
)

// Validator Структура дополнительной проверки таймеров переписки и работы менеджеров
type Validator struct {
	Limit            int64
	GetDb            func() *sql.DB
	stop             chan bool
	getTableName     func(string) string
	sendCronMessages func(rooms.Room, []messages.Message)
}

// New Сбор валидатора с интервалом
func New(limit int64, getDb func() *sql.DB) *Validator {
	return &Validator{
		Limit: limit,
		GetDb: getDb,
		stop:  make(chan bool),
	}
}

// Check Функция таймера
func (l *Validator) Check() {
	ticker := time.NewTicker(time.Duration(30) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				l.Run()
			case <-l.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

// SetTableName Функция возвращает название таблицы
func (l *Validator) SetTableName(f func(string) string) {
	l.getTableName = f
}

// SetSendCronMessages Функция для отправки рассылки
func (l *Validator) SetSendCronMessages(f func(rooms.Room, []messages.Message)) {
	l.sendCronMessages = f
}

// GetTableName Функция возвращает название таблицы
func (l *Validator) GetTableName(name string) string {
	if l.getTableName != nil {
		return l.getTableName(name)
	}
	return "rooms"
}

// Run Функция выполняющая проверку
func (l *Validator) Run() {
	// И так для начала собераем все открытые комнаты
	chats, _ := rooms.GetAllIsOpen(l.GetTableName("rooms"), l.GetDb())
	if len(chats) == 0 {
		return
	}
	// Далее нужно отфильтровать комнаты у которых срок ответа истек
	chats = l.FilterRooms(chats)
	// Далее нужно разделить комнаты на пустые, с группой и менеджером
	e, g, m := l.RangeRoom(chats)
	// Далее нужно что то сделать с путсыми комнатами
	l.ExecEmptyRooms(e)
	// Далее нужно переслать все в другую группу
	l.ExecGroupsRooms(g)
	// Далее нужно переслать все в другую группу
	l.ExecManagersRooms(m)
}

// FilterRooms отфильтровывает комнаты у которых срок ответа истек
func (l *Validator) FilterRooms(r []rooms.Room) []rooms.Room {
	for i, v := range r {
		m, err := messages.GetLastMessage(v.Room, l.GetTableName("message"), l.GetDb())
		if err != nil {
			// Пока оставляю комнату на обработку, так как у комнеты нет сообщений
			fmt.Println(err, v.Room)
			continue
		}
		t := time.Now()
		d, err := time.Parse("2006-01-02 15:04:05", m.Datetime)
		if err != nil {
			fmt.Println(err)
			continue
		}
		n := d.Sub(t)

		fmt.Println("now", t)
		fmt.Println("ago", d)
		fmt.Println("Message", m.ID, v.Room, m.Room)
		if n.Hours() > 0 {
			// Более часа ждать нельзя
			continue
		}
		if int64(n.Seconds()) > l.Limit {
			// Время вышло
			continue
		}
		fmt.Println("Нужно удалить", m.ID, r)
		copy(r[i:], r[i+1:])
		r[len(r)-1] = rooms.Room{}
		r = r[:len(r)-1]
		fmt.Println("Нужно удалить", r)
	}
	return r
}

// RangeRoom Функция разделяет массив по категориям
func (l *Validator) RangeRoom(r []rooms.Room) ([]rooms.Room, []rooms.Room, []rooms.Room) {
	e := []rooms.Room{}
	g := []rooms.Room{}
	m := []rooms.Room{}
	for _, v := range r {
		if v.ChatID == 0 && v.GroupID == 0 {
			e = append(e, v)
		}
		if v.ChatID == 0 && v.GroupID > 0 {
			g = append(g, v)
		}
		if v.ChatID > 0 {
			m = append(m, v)
		}
	}
	return e, g, m
}

// ExecEmptyRooms Комната пуста. Что с ней делать?
func (l *Validator) ExecEmptyRooms(r []rooms.Room) {}

// ExecGroupsRooms Комната групп меняется если есть дочерняя группа
func (l *Validator) ExecGroupsRooms(r []rooms.Room) {
	for i, value := range r {
		// Ищем дочернюю группу
		g, err := groups.GetParents(value.GroupID, l.GetTableName("groups"), l.GetDb())
		if err != nil || len(g) == 0 {
			fmt.Println("Find Groups error", err)
			continue
		}
		// Далее устанавливаем новую группу
		value.GroupID = g[0].ID
		err = rooms.Update(value, l.GetTableName("rooms"), l.GetDb())
		if err != nil {
			fmt.Println("ExecGroupsRooms Udate Group error", err)
			continue
		}
		r[i] = value
		// Далее делаем рассылку в новую группу
		msg, err := rooms.GetMessagesRoom(value.ID, l.GetTableName("message"), l.GetTableName("rooms"), l.GetDb())
		if err != nil {
			continue
		}
		if len(msg) > 0 {
			l.SendCronMessages(value, msg)
			// Далее обновляем время последнего сообщения для смещения
			m := messages.Message{}
			for _, u := range msg {
				if m.ID < u.ID {
					m = u
				}
			}
			m.GroupID = value.GroupID
			messages.Update(m, l.GetTableName("message"), l.GetDb())
		}
	}
}

// ExecManagersRooms Комната выпинывает менеджера и делает рассылку в ту же группу
func (l *Validator) ExecManagersRooms(r []rooms.Room) {
	for _, value := range r {
		fmt.Println("ExecManagersRooms", value)
	}
}

// SendCronMessages рассылка в телеграмм группу
func (l *Validator) SendCronMessages(room rooms.Room, msgs []messages.Message) {
	if l.sendCronMessages != nil {
		l.sendCronMessages(room, msgs)
	}
}
