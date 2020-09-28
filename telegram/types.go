package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// BotApp Основная структура приложения бота
type BotApp struct {
	API                   *tgbotapi.BotAPI        // API телеграмма
	Updates               tgbotapi.UpdatesChannel // Канал обновлений
	ActiveContactRequests []int64
	ActivePostRequests    []PostRequests
	findUser              func(chatID int64) User
	createUser            func(chatID int64) User
	updateUser            func(user User)
	onMessage             func(chatID int64, message string)
	onComands             func(chatID int64, message PostRequests)
}

// PostRequests Структура техника
type PostRequests struct {
	ChatID  int64
	Message string
	RoomID  int64
}

// User Структура пользователя
type User struct {
	ID        int64  `json:"id"`
	ChatID    int64  `json:"chatID"`
	UserID    int64  `json:"userID"`
	Login     string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Status    bool   `json:"status"`
	Date      string `json:"date"`
	Reghash   string `json:"regHash"`
}
