package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Start Основной цикл бота
func (app *BotApp) Start() {
	app.Info("Запуск сервиса телеграм-бота")
	for update := range app.Updates {
		if update.Message != nil {
			// Если сообщение есть и его длина больше 0 -> начинаем обработку
			app.analyzeUpdate(update)
		}
	}
}

// FindUser Промежуточная функция поиска менеджеров
func (app *BotApp) FindUser(chatID int64) User {
	return app.findUser(chatID)
}

// CreateUser Промежуточная функция поиска менеджеров
func (app *BotApp) CreateUser(chatID int64) User {
	return app.createUser(chatID)
}

// UpdateUser Промежуточная функция поиска менеджеров
func (app *BotApp) UpdateUser(user User) {
	app.updateUser(user)
}

// analyzeUpdate Начало обработки сообщения
func (app *BotApp) analyzeUpdate(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	user := app.FindUser(chatID)
	if user.ChatID == 0 {
		app.Info("%v", update)
		user = app.CreateUser(chatID)
	}
	app.analyzeUser(update, user)
}

func (app *BotApp) findContactRequestID(chatID int64) bool {
	for _, v := range app.ActiveContactRequests {
		if v == chatID {
			return true
		}
	}
	return false
}

func (app *BotApp) addContactRequestID(chatID int64) {
	app.ActiveContactRequests = append(app.ActiveContactRequests, chatID)
}

func (app *BotApp) deleteContactRequestID(chatID int64) {
	for i, v := range app.ActiveContactRequests {
		if v == chatID {
			copy(app.ActiveContactRequests[i:], app.ActiveContactRequests[i+1:])
			app.ActiveContactRequests[len(app.ActiveContactRequests)-1] = 0
			app.ActiveContactRequests = app.ActiveContactRequests[:len(app.ActiveContactRequests)-1]

		}
	}
}
func (app *BotApp) addPostRequests(data PostRequests) {
	app.ActivePostRequests = append(app.ActivePostRequests, data)
}

func (app *BotApp) deletePostRequests(chatID int64) {
	for i, v := range app.ActivePostRequests {
		if v.ChatID == chatID {
			copy(app.ActivePostRequests[i:], app.ActivePostRequests[i+1:])
			app.ActivePostRequests[len(app.ActivePostRequests)-1] = PostRequests{}
			app.ActivePostRequests = app.ActivePostRequests[:len(app.ActivePostRequests)-1]
		}
	}
}
func (app *BotApp) findPostRequests(chatID int64) (PostRequests, bool) {
	for _, v := range app.ActivePostRequests {
		if v.ChatID == chatID {
			return v, true
		}
	}
	return PostRequests{}, false
}

// Проверка принятого контакта
func (app *BotApp) checkRequestContactReply(update tgbotapi.Update) {
	if update.Message.Contact != nil { // Проверяем, содержит ли сообщение контакт
		from := update.Message.From
		if update.Message.Contact.UserID == from.ID { // Проверяем действительно ли это контакт отправителя
			user := app.FindUser(update.Message.Chat.ID)
			if user.ChatID > 0 {
				user.FirstName = from.FirstName
				user.LastName = from.LastName
				user.Login = from.UserName
				app.UpdateUser(user)
			} else {
				app.CreateUser(update.Message.Chat.ID)
				app.Send(update.Message.Chat.ID, "Произошла ошибка! Попробуйте позже")
				return
			}
			app.deleteContactRequestID(update.Message.Chat.ID) // Удаляем ChatID из списка ожидания
			app.Send(update.Message.Chat.ID, fmt.Sprintf("Ваш код %s для подтверждения регистрации!", user.Reghash))

		} else {
			app.Send(update.Message.Chat.ID, "Профиль, который вы предоставили, принадлежит не вам!")
			app.RequestContact(update.Message.Chat.ID)
		}
	} else {
		app.Send(update.Message.Chat.ID, "Если вы не предоставите ваш профиль, вы не сможете пользоваться системой!")
		app.RequestContact(update.Message.Chat.ID)
	}
}

// RequestContact Запросить номер телефона
func (app *BotApp) RequestContact(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Согласны ли вы предоставить ваш профиль для регистрации в системе?")
	acceptButton := tgbotapi.NewKeyboardButtonContact("Да")
	declineButton := tgbotapi.NewKeyboardButton("Нет")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{acceptButton, declineButton})
	app.API.Send(msg)
	app.addContactRequestID(chatID)
}
