package telegram

import (
	"fmt"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *BotApp) analyzeUser(update tgbotapi.Update, user User) {
	// Получаем id чата чтоб бот понимал кому слать ответ
	chatID := update.Message.Chat.ID
	// Проверяем ожидаем ли ответа от менеджера
	if post, ok := app.findPostRequests(chatID); ok {
		re := regexp.MustCompile(`([0-9]+)`)
		validID := re.FindAllString(update.Message.Text, -1)
		app.deletePostRequests(chatID)
		if len(validID) > 0 {
			if roomID, err := strconv.Atoi(validID[0]); err == nil {
				post.RoomID = int64(roomID)
				app.onComands(chatID, post)
			} else {
				app.SendOpenButton(chatID, "Не удалось распознать номер чата")
				return
			}
		} else {
			app.SendOpenButton(chatID, "Не удалось распознать номер чата")
			return
		}
		return
	}
	// Проверяем менеджер ли это или не менеджер
	if user.UserID > 0 && user.Status == true {
		switch true {
		// Прикол отсебячина
		case update.Message.Text == "ты дурак":
			app.Send(chatID, "Сам дурак")
			return
		case update.Message.Text == "/start":
			app.SendOpenButton(chatID, "Какую информацию вы хотите получить?")
			return
		default:
			app.onMessage(chatID, update.Message.Text)
			return
		}
	} else if user.UserID > 0 && user.Status == false {
		app.Send(chatID, "Вам не предоставлен доступ к сервису.")
		return
	} else if len(user.Login) > 0 && len(user.Reghash) > 0 {
		app.Send(chatID, fmt.Sprintf("Ваш код %s для подтверждения регистрации!", user.Reghash))
		return
	} else {
		// Если номера нет, то проверяем ждём ли мы контакт от этого ChatID
		if app.findContactRequestID(chatID) {
			app.checkRequestContactReply(update) // Если да -> проверяем
			return
		}
		app.RequestContact(chatID) // Если нет -> запрашиваем его
	}
}

// ReloadButton Кнопки по умолчанию
func (app *BotApp) ReloadButton(msg *tgbotapi.MessageConfig) {
	markup := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Занять чат с клиентом"),
			tgbotapi.NewKeyboardButton("Получить переписку чата"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Выйти из чата"),
			tgbotapi.NewKeyboardButton("Закрыть чат"),
		),
	)
	msg.ReplyMarkup = &markup
}

// RoomsButtonSend Кнопки по умолчанию
func (app *BotApp) RoomsButtonSend(chatID int64, rooms []int64, message string) {
	app.addPostRequests(PostRequests{chatID, message, 0})
	msg := tgbotapi.NewMessage(chatID, message)
	allkeys := [][]tgbotapi.KeyboardButton{}
	keys := []tgbotapi.KeyboardButton{}
	for k, v := range rooms {
		if k%4 == 3 {
			allkeys = append(allkeys, tgbotapi.NewKeyboardButtonRow(keys...))
			keys = []tgbotapi.KeyboardButton{}
		}
		keys = append(keys, tgbotapi.NewKeyboardButton(fmt.Sprintf("Чат номер %d", v)))
	}
	if len(keys) > 0 {
		allkeys = append(allkeys, keys)
	}
	markup := tgbotapi.NewReplyKeyboard(allkeys...)
	msg.ReplyMarkup = &markup
	app.API.Send(msg)
}

// Send Отправка сообщения пользователю телеграмм
func (app *BotApp) Send(chatID int64, msg string) {
	option := tgbotapi.NewMessage(chatID, msg)
	option.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	option.ParseMode = tgbotapi.ModeMarkdown
	app.API.Send(option)
}

// SendOpenButton Отправка сообщения пользователю телеграмм с кнопками
func (app *BotApp) SendOpenButton(chatID int64, msg string) {
	option := tgbotapi.NewMessage(chatID, msg)
	app.ReloadButton(&option)
	tgbotapi.NewRemoveKeyboard(false)
	option.ParseMode = tgbotapi.ModeMarkdown
	app.API.Send(option)
}
