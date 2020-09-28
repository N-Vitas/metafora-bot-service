package webapp

import "fmt"

// Message Структура сообщения
type Message struct {
	Author string      `json:"author"`
	Body   string      `json:"body"`
	Action string      `json:"action"`
	ChatID int64       `json:"chatID"`
	Params interface{} `json:"params"`
}

// ClientChatMessage Структура сообщения клиента
type ClientChatMessage struct {
	ID       string   `json:"id"`
	Action   string   `json:"action"`
	Img      string   `json:"img"`
	ChatID   int64    `json:"chatID"`
	GroupID  int64    `json:"groupID"`
	ReplicID int64    `json:"replicID"`
	ChatRoom string   `json:"room"`
	Message  string   `json:"message"`
	Date     string   `json:"datetime"`
	Status   int64    `json:"status"`
	Type     string   `json:"type"`
	DataType []string `json:"dataType"`
}

const (
	// InfoColor Синий цвет
	InfoColor = "\033[1;34m%v\033[0m"
	// NoticeColor Голубой цвет
	NoticeColor = "\033[1;36m%v\033[0m"
	// WarningColor Желтый цвет
	WarningColor = "\033[1;33m%v\033[0m"
	// ErrorColor Красный цвет
	ErrorColor = "\033[1;31m%v\033[0m"
	// DebugColor Серый цвет
	DebugColor = "\033[0;36m%v\033[0m"
)

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
