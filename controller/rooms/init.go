package rooms

import (
	"database/sql"
	"fmt"
	"metafora-bot-service/controller/messages"
)

// Room Структура комнаты
type Room struct {
	ID          int64  `json:"id"`
	ChatID      int64  `json:"chatID"`
	GroupID     int64  `json:"groupID"`
	ReplicID    int64  `json:"replicID"`
	LastMessage int64  `json:"lastmessage"`
	MessagesID  string `json:"messagesID"`
	Room        string `json:"chatRoom"`
	Mute        bool   `json:"mute"`
	Status      int64  `json:"status"`
	Datetime    string `json:"date"`
}

// Init Создание базы данных для сообщений
func Init(table string, db *sql.DB) error {
	// query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	// 	id int(9) NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	chatID int(9) DEFAULT 0,
	// 	groupID int(9) DEFAULT 0,
	// 	replicID int(9) DEFAULT 0,
	// 	lastmessage int(9) NOT NULL DEFAULT 0,
	// 	messagesID text COLLATE utf8mb4_unicode_ci DEFAULT '',
	// 	chatRoom varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	// 	date datetime NOT NULL DEFAULT current_timestamp(),
	// 	mute tinyint(1) NOT NULL DEFAULT 0,
	// 	status int(5) NOT NULL DEFAULT 1,
	// 	UNIQUE KEY chatRoom (chatRoom)
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	// 	table,
	// )
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		chatID integer DEFAULT 0,
		groupID integer DEFAULT 0,
		replicID integer DEFAULT 0,
		lastmessage integer NOT NULL DEFAULT 0,
		messagesID TEXT DEFAULT '',
		chatRoom TEXT NOT NULL DEFAULT '' UNIQUE,
		date TEXT NOT NULL DEFAULT '',
		mute integer NOT NULL DEFAULT 0,
		status integer NOT NULL DEFAULT 1
	  );`,
		table,
	)
	_, err := db.Exec(query)
	return err
}

// Get Получение комнаты
func Get(room string, table string, db *sql.DB) (Room, error) {
	s := Room{Room: room}
	query := fmt.Sprintf(`SELECT id, chatID, groupID, replicID, lastmessage, messagesID, date, mute, status FROM %s WHERE chatRoom = '%s' AND status > 0`, table, room)
	err := db.QueryRow(query).Scan(&s.ID, &s.ChatID, &s.GroupID, &s.ReplicID, &s.LastMessage, &s.MessagesID, &s.Datetime, &s.Mute, &s.Status)
	return s, err
}

// GetByID Получение комнаты
func GetByID(id int64, table string, db *sql.DB) (Room, error) {
	s := Room{ID: id}
	query := fmt.Sprintf(`SELECT chatRoom, chatID, groupID, replicID, lastmessage, messagesID, date, mute, status FROM %s WHERE id = %d`, table, id)
	err := db.QueryRow(query).Scan(&s.Room, &s.ChatID, &s.GroupID, &s.ReplicID, &s.LastMessage, &s.MessagesID, &s.Datetime, &s.Mute, &s.Status)
	return s, err
}

// GetAllIsOpen Получение всех сообщений комнаты по ID
func GetAllIsOpen(table string, db *sql.DB) ([]Room, error) {
	m := []Room{}
	ID := sql.NullInt64{}
	ChatID := sql.NullInt64{}
	GroupID := sql.NullInt64{}
	ReplicID := sql.NullInt64{}
	LastMessage := sql.NullInt64{}
	MessagesID := sql.NullString{}
	ChatRoom := sql.NullString{}
	Mute := sql.NullBool{}
	Status := sql.NullInt64{}
	Datetime := sql.NullString{}
	query := fmt.Sprintf("SELECT id, chatID, groupID, replicID, lastmessage, messagesID, chatRoom, date, mute, status FROM %s WHERE status > 0", table)
	rows, err := db.Query(query)
	if err != nil {
		return m, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&ID, &ChatID, &GroupID, &ReplicID, &LastMessage, &MessagesID, &ChatRoom, &Datetime, &Mute, &Status)
		if err != nil {
			continue
		}
		m = append(m, Room{
			ID:          ID.Int64,
			ChatID:      ChatID.Int64,
			GroupID:     GroupID.Int64,
			ReplicID:    ReplicID.Int64,
			LastMessage: LastMessage.Int64,
			MessagesID:  MessagesID.String,
			Room:        ChatRoom.String,
			Mute:        Mute.Bool,
			Status:      Status.Int64,
			Datetime:    Datetime.String,
		})
	}
	return m, err
}

// FindByManager Получение комнаты с менеджером
func FindByManager(chatID int64, table string, db *sql.DB) (Room, error) {
	s := Room{ChatID: chatID}
	query := fmt.Sprintf(`SELECT id, chatRoom, groupID, replicID, lastmessage, messagesID, date, mute, status FROM %s WHERE chatID = %d AND mute = 0`, table, chatID)
	err := db.QueryRow(query).Scan(&s.ID, &s.Room, &s.GroupID, &s.ReplicID, &s.LastMessage, &s.MessagesID, &s.Datetime, &s.Mute, &s.Status)
	return s, err
}

// Create Новая комната
func Create(room string, table string, db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`INSERT INTO %s (chatRoom, replicID) VALUES ('%s', 1)`, table, room))
	return err
}

// Update Обновление комнаты
func Update(room Room, table string, db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`UPDATE %s SET chatID=%d, groupID=%d, replicID=%d, lastmessage=%d, messagesID='%s', mute=%v, status=%d WHERE chatRoom = '%s'`,
		table, room.ChatID, room.GroupID, room.ReplicID, room.LastMessage, room.MessagesID, room.Mute, room.Status, room.Room))
	return err
}

// GoRoomMagager Менеджер занимает комнату
func GoRoomMagager(chatID int64, roomID int64, table string, db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`UPDATE %s SET chatID=%d, mute=0 WHERE id = %d`,
		table, chatID, roomID))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf(`UPDATE %s SET mute=1 WHERE chatID = %d AND id != %d`,
		table, chatID, roomID))
	return err
}

// ExitRoomMagager Менеджер освобождает комнату
func ExitRoomMagager(chatID int64, table string, db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`UPDATE %s SET chatID = 0, mute=0 WHERE chatID = %d and mute=0`, table, chatID))
	return err
}

// ClosedRoom Менеджер закрывает комнату
func ClosedRoom(roomID int64, table string, db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`UPDATE %s SET status = 0, mute=1 WHERE id = %d`, table, roomID))
	return err
}

// GetMessagesRoom Получение всех сообщений комнаты по ID
func GetMessagesRoom(roomID int64, tableMessage string, tableRoom string, db *sql.DB) ([]messages.Message, error) {
	m := []messages.Message{}
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
	query := fmt.Sprintf(`SELECT m.id, m.img, m.chatRoom, m.message, m.chatID, m.groupID, m.replicID, m.status, m.date, m.type, m.dataType 
	FROM %s r INNER JOIN %s m ON m.chatRoom = r.chatRoom WHERE r.id = %d`, tableRoom, tableMessage, roomID)
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
		m = append(m, messages.Message{
			ID:       ID.Int64,
			Img:      Img.String,
			Room:     Room.String,
			Message:  Msge.String,
			ChatID:   ChatID.Int64,
			GroupID:  GroupID.Int64,
			ReplicID: ReplicID.Int64,
			Status:   Status.Int64,
			Datetime: Datetime.String,
			Type:     Type.String,
			DataType: DataType.String,
		})
	}
	return m, err
}
