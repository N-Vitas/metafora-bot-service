package controller

import (
	"fmt"
	"metafora-bot-service/controller/crontimes"
	"metafora-bot-service/controller/groupmanager"
	"metafora-bot-service/controller/groups"
	"metafora-bot-service/controller/managers"
	"metafora-bot-service/controller/messages"
	"metafora-bot-service/controller/questions"
	"metafora-bot-service/controller/rooms"
	"metafora-bot-service/controller/settings"
	"metafora-bot-service/controller/user"
	"metafora-bot-service/database"
	"metafora-bot-service/telegram"
	"strings"
)

// Controller Структура логики чат бота
type Controller struct {
	DB               *database.SessionDb
	Host             string
	FolderID         string
	BotToken         string
	Tables           []string
	DurationClients  int64
	bot              *telegram.BotApp
	onMessage        func(chatRoom string, chatID int64, message string, params interface{})
	deleteRoomClient func(chatRoom string, chatID int64, message string, params interface{})
	Cron             *crontimes.Validator
}

// Init Инициализация контроллера
func Init(db *database.SessionDb) *Controller {
	// findUser func(chatID int64) User, createUser func(chatID int64) User
	c := &Controller{
		DB:       db,
		Host:     "0.0.0.0:8082",
		FolderID: "1MmY3uvmNcOk1Zg_9gaW2KMd7qly2H2WQ",
		Tables: []string{
			"wp_chatbottelegram_bot_tell",
			"wp_chatbottelegram_clients",
			"wp_chatbottelegram_groups",
			"wp_chatbottelegram_managers",
			"wp_chatbottelegram_group_managers",
			"wp_chatbottelegram_messages",
			"wp_chatbottelegram_rooms",
			"wp_chatbottelegram_settings",
			"wp_chatbottelegram_user",
		},
		DurationClients: 30,
		Cron:            crontimes.New(60, db.GetDb),
	}
	if err := messages.Init(c.GetTableName("message"), db.GetDb()); err != nil {
		Error("Messages Init error : %v", err)
	}
	if err := groups.Init(c.GetTableName("groups"), db.GetDb()); err != nil {
		Error("Groups Init error : %v", err)
	}
	if err := managers.Init(c.GetTableName("managers"), db.GetDb()); err != nil {
		Error("Managers Init error : %v", err)
	}
	if err := rooms.Init(c.GetTableName("rooms"), db.GetDb()); err != nil {
		Error("Rooms Init error : %v", err)
	}
	if err := questions.Init(c.GetTableName("bot"), db.GetDb()); err != nil {
		Error("Questions Init error : %v", err)
	}
	if err := settings.Init(c.GetTableName("settings"), db.GetDb()); err != nil {
		Error("Settings Init error : %v", err)
	}
	if err := groupmanager.Init(c.GetTableName("group_manager"), db.GetDb(), c.GetTableName("groups"), c.GetTableName("managers")); err != nil {
		Error("Group-Managers Init error : %v", err)
	}
	if err := user.Init(c.GetTableName("user"), db.GetDb()); err != nil {
		Error("User Init error : %v", err)
	}
	// Берем настройки в базе
	if s, e := settings.Get(c.GetTableName("settings"), c.DB.GetDb()); e == nil {
		Notice("Время ожидания реплики бота в секундах %d", s.DurationClients)
		c.DurationClients = s.DurationClients
		c.Host = s.HostService
		c.FolderID = s.GoogleFolder
		c.BotToken = s.Token
		c.Cron.Limit = s.DurationManagers
	}
	c.bot = telegram.NewTelegramApp(c.BotToken, c.FindUser, c.CreateUser, c.UpdateUser, c.OnMessage, c.OnComands)
	go c.bot.Start()
	c.Cron.SetTableName(c.GetTableName)
	c.Cron.SetSendCronMessages(c.SendCronMessages)
	c.Cron.Check()
	return c
}

// CreateUser Передача создания менеджера
func (c *Controller) CreateUser(chatID int64) telegram.User {
	u := telegram.User{ChatID: chatID}
	id, err := managers.NewUser(c.GetTableName("managers"), c.DB.GetDb(), chatID)
	if err != nil {
		Error("Ошибка создания менеджера : %v", err)
	}
	u.ID = id
	return u
}

// FindUser Передача поиска менеджера по ID чата
func (c *Controller) FindUser(chatID int64) telegram.User {
	u1 := telegram.User{}
	u2, err := managers.FindManagerByChatID(chatID, c.GetTableName("managers"), c.DB.GetDb())
	if err != nil {
		Error("Ошибка поиска менеджера : %v", err)
		return u1
	}
	u1.ID = u2.ID
	u1.ChatID = u2.ChatID
	u1.UserID = u2.UserID
	u1.Login = u2.UserName
	u1.FirstName = u2.FirstName
	u1.LastName = u2.LastName
	u1.Status = u2.Status > 0
	u1.Date = u2.Datetime
	u1.Reghash = u2.Reghash
	return u1
}

// UpdateUser Передача поиска менеджера по ID чата
func (c *Controller) UpdateUser(user telegram.User) {
	manager, err := managers.FindManagerByChatID(user.ChatID, c.GetTableName("managers"), c.DB.GetDb())
	if err != nil {
		Error("Ошибка поиска менеджера : %v", err)
	}
	manager.ID = user.ID
	manager.UserID = user.UserID
	manager.UserName = user.Login
	manager.FirstName = user.FirstName
	manager.LastName = user.LastName
	err = managers.Update(manager, c.GetTableName("managers"), c.DB.GetDb())
	if err != nil {
		Error("Ошибка обновления менеджера : %v", err)
	}
}

// GetTableName Получение полного имени таблицы
func (c *Controller) GetTableName(name string) string {
	for _, table := range c.Tables {
		if strings.Index(table, name) != -1 {
			return table
		}
	}
	return name
}

// Info Логирование в консоль синего цвета
func Info(template string, val ...interface{}) {
	fmt.Printf("\033[1;34m"+template+"\033[0m\n", val...)
}

// Notice Логирование в консоль голубого цвета
func Notice(template string, val ...interface{}) {
	fmt.Printf("\033[1;36m"+template+"\033[0m\n", val...)
}

// Warning Логирование в консоль желтого цвета
func Warning(template string, val ...interface{}) {
	fmt.Printf("\033[1;33m"+template+"\033[0m\n", val...)
}

// Error Логирование в консоль красного цвета
func Error(template string, val ...interface{}) {
	fmt.Printf("\033[1;31m"+template+"\033[0m\n", val...)
}

// Debug Логирование в консоль серого цвета
func Debug(template string, val ...interface{}) {
	fmt.Printf("\033[0;36m"+template+"\033[0m\n", val...)
}

// GetSettings Запрос настроек
func (c *Controller) GetSettings() settings.Settings {
	s, e := settings.Get(c.GetTableName("settings"), c.DB.GetDb())
	if e != nil {
		Error("Ошибка запроса настроек %s", e.Error())
		return settings.Settings{}
	}
	return s
}

// MaxSort Проверка последней реплики
func (c *Controller) MaxSort() int64 {
	return questions.MaxSort(c.GetTableName("bot"), c.DB.GetDb())
}

// GetAllReplics Получение всех реплик
func (c *Controller) GetAllReplics() []questions.Replic {
	return questions.GetAll(c.GetTableName("bot"), c.DB.GetDb())
}
