package telegram

import (
	"fmt"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	// BotAppAPIKey Ключь основного телеграмм бота
	BotAppAPIKey = "573830185:AAEkh87CSsIZvWJVGksVX1vxUDGqcFF3dUM"
	// BotAppDebugAPIKey Ключь тестового телеграмм бота
	BotAppDebugAPIKey = "295029059:AAHefWQCE1yCbnZwjLUWkjHa7haVrQlH4Ao"
	// BotAppUpdateOffset видимо время задержки
	BotAppUpdateOffset = 0
	// BotAppUpdateTimeout время ожидания обновления
	BotAppUpdateTimeout = 64
)

// NewTelegramApp Создание приложения бота
func NewTelegramApp(token string, findUser func(chatID int64) User, createUser func(chatID int64) User, updateUser func(User),
	onMessage func(chatID int64, message string), onComands func(chatID int64, message PostRequests)) *BotApp {
	app := &BotApp{
		findUser:   findUser,
		createUser: createUser,
		updateUser: updateUser,
		onMessage:  onMessage,
		onComands:  onComands,
	}
	app.Init(token)
	return app
}

// Init Инициализация бота
func (app *BotApp) Init(token string) {
	botAPI, err := tgbotapi.NewBotAPI(app.getToken(token)) // Инициализация API
	if err != nil {
		app.Info("Ошибка инициализации апи телеграм-бота %s %s", err.Error(), token)
		return
	}
	app.API = botAPI
	botUpdate := tgbotapi.NewUpdate(BotAppUpdateOffset) // Инициализация канала обновлений
	botUpdate.Timeout = BotAppUpdateTimeout
	botUpdates, err := app.API.GetUpdatesChan(botUpdate)
	if err != nil {
		app.Info("Ошибка подключения тригера обновления %s", err.Error())
		return
	}
	app.Updates = botUpdates
}

func (app *BotApp) getToken(token string) string {
	if ok, _ := regexp.MatchString(`[0-9]{9}:[0-9a-zA-Z]{35}`, token); ok {
		return token
	}
	return BotAppAPIKey
}

// Info Вывод информации в консоль
func (app *BotApp) Info(template string, values ...interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf(t+" \033[1;32m[telegram][info]\033[0m \033[1;33m"+template+"\033[0m\n", values...)
}
