package messages

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Message Структура сообщения
type Message struct {
	ID       int64  `json:"id"`
	Img      string `json:"img"`
	Room     string `json:"chatRoom"`
	Message  string `json:"message"`
	ChatID   int64  `json:"chatID"`
	GroupID  int64  `json:"groupID"`
	ReplicID int64  `json:"replicID"`
	Status   int64  `json:"status"`
	Datetime string `json:"datetime"`
	Type     string `json:"type"`
	DataType string `json:"dataType"`
}

// Init Создание базы данных для сообщений
func Init(table string, db *sql.DB) error {
	// query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	// 	id int(9) NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	img text COLLATE utf8mb4_unicode_ci DEFAULT '',
	// 	chatID int(9) DEFAULT 0,
	// 	groupID int(9) DEFAULT 0,
	// 	replicID int(9) DEFAULT 0,
	// 	chatRoom varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	// 	message text COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	// 	date datetime NOT NULL DEFAULT current_timestamp(),
	// 	status int(5) NOT NULL DEFAULT 1,
	// 	type varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT 'text',
	// 	dataType varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '[]'
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	// 	table,
	// )
	query := `CREATE TABLE IF NOT EXISTS ` + table + ` (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		img text DEFAULT '',
		chatID integer DEFAULT 0,
		groupID integer DEFAULT 0,
		replicID integer DEFAULT 0,
		chatRoom text NOT NULL DEFAULT '',
		message text NOT NULL DEFAULT '',
		date TEXT NOT NULL DEFAULT '',
		status integer NOT NULL DEFAULT 1,
		type text DEFAULT 'text',
		dataType text DEFAULT ''
	  );`
	_, err := db.Exec(query)
	return err
}

// GetAllInID Получение всех сообщений комнаты по ID
func GetAllInID(messagesID string, table string, db *sql.DB) ([]Message, error) {
	m := []Message{}
	re := regexp.MustCompile(`([0-9]+)`)
	validID := re.FindAllString(messagesID, -1)
	if len(validID) > 0 {
		ID := sql.NullInt64{}
		Img := sql.NullString{}
		Room := sql.NullString{}
		Msge := sql.NullString{}
		ChatID := sql.NullInt64{}
		GroupID := sql.NullInt64{}
		ReplicID := sql.NullInt64{}
		Status := sql.NullInt64{}
		Datetime := sql.NullString{}
		Type := sql.NullString{}
		DataType := sql.NullString{}
		query := fmt.Sprintf("SELECT id, img, chatRoom, message, chatID, groupID, replicID, status, date, type, dataType FROM %s WHERE id in(%s)", table, strings.Join(validID, ","))
		rows, err := db.Query(query)
		if err != nil {
			return m, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&ID, &Img, &Room, &Msge, &ChatID, &GroupID, &ReplicID, &Status, &Datetime, &Type, &DataType)
			if err != nil {
				continue
			}
			m = append(m, Message{
				ID.Int64,
				Img.String,
				Room.String,
				Msge.String,
				ChatID.Int64,
				GroupID.Int64,
				ReplicID.Int64,
				Status.Int64,
				Datetime.String,
				Type.String,
				DataType.String,
			})
		}
		return m, err
	}
	return m, errors.New("Строка не содержит ID сообщений")
}

// NewMessage Создание нового сообщения
func NewMessage(table string, db *sql.DB, img, chatRoom, message, types, dataType string, chatID, groupID, replicID int64) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (img, chatID, groupID, replicID, chatRoom, message, date, type, dataType) 
		VALUES ('%s', %d, %d, %d, '%s', '%s', DATETIMES, '%s', '%s')`,
		table, img, chatID, groupID, replicID, chatRoom, message, types, dataType)
	query = strings.Replace(query, "DATETIMES", "strftime('%Y-%m-%d %H:%M:%S','now')", -1)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
