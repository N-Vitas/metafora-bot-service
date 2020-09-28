package managers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Manager Структура менеджера
type Manager struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userID"`
	ChatID    int64  `json:"chatID"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
	Reghash   string `json:"reghash"`
	Status    int64  `json:"status"`
	Datetime  string `json:"date"`
}

// Init Создание базы данных для сообщений
func Init(table string, db *sql.DB) error {
	// query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	// 	id int(9) NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	userID bigint(20) DEFAULT 0,
	// 	chatID int(9) DEFAULT 0,
	// 	firstname varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
	// 	lastname varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
	// 	username varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
	// 	reghash varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
	// 	date datetime NOT NULL DEFAULT current_timestamp(),
	// 	status int(5) NOT NULL DEFAULT 1
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	// 	table,
	// )
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		userID integer DEFAULT 0,
		chatID integer DEFAULT 0,
		firstname text DEFAULT '',
		lastname text DEFAULT '',
		username text DEFAULT '',
		reghash text DEFAULT '',
		date text NOT NULL DEFAULT '',
		status integer NOT NULL DEFAULT 1
	  );`,
		table,
	)
	_, err := db.Exec(query)
	return err
}

// FindManagerByChatID Поиск менеджера по ID чата
func FindManagerByChatID(chatID int64, table string, db *sql.DB) (Manager, error) {
	s := Manager{ChatID: chatID}
	query := fmt.Sprintf(`SELECT id, userID, firstname, lastname, username, reghash, date, status FROM %s WHERE chatID = %d`, table, chatID)
	err := db.QueryRow(query).Scan(&s.ID, &s.UserID, &s.FirstName, &s.LastName, &s.UserName, &s.Reghash, &s.Datetime, &s.Status)
	return s, err
}

// Get Получение менеджера
func Get(id int64, table string, db *sql.DB) (Manager, error) {
	s := Manager{ID: id}
	query := fmt.Sprintf(`SELECT userID, chatID, firstname, lastname, username, reghash, date, status FROM %s WHERE id = %d`, table, id)
	err := db.QueryRow(query).Scan(&s.UserID, &s.ChatID, &s.FirstName, &s.LastName, &s.UserName, &s.Datetime, &s.Status)
	return s, err
}

// NewUser Создание нового менеджера
func NewUser(table string, db *sql.DB, chatID int64) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (chatID, reghash, date) VALUES (%d, '%s', DATETIMES)`, table, chatID, generate())
	query = strings.Replace(query, "DATETIMES", "strftime('%Y-%m-%d %H:%M:%S','now')", -1)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Update Новая комната
func Update(user Manager, table string, db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`UPDATE %s SET userID=%d, firstname='%s', lastname='%s', username='%s', reghash='%s', status=%d WHERE chatID=%d`,
		table, user.UserID, user.FirstName, user.LastName, user.UserName, user.Reghash, user.Status, user.ChatID))
	return err
}
func generate() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	length := 6
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "ExcbsVQs"
	return str
}
