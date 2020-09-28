package user

import (
	"database/sql"
	"fmt"
)

// User Структура пользователя
type User struct {
	ID      int64  `json:"id"`
	Login   string `json:"login"`
	Pasword string `json:"password"`
	Name    string `json:"name"`
	Type    string `json:"userType"`
	Blocked int64  `json:"blocked"`
}

// Init Создание базы данных для сообщений
func Init(table string, db *sql.DB) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		login text DEFAULT '',
		password text DEFAULT '',
		name text DEFAULT '',
		userType text DEFAULT '',
		blocked integer NOT NULL DEFAULT 1
	  );`,
		table,
	)
	_, err := db.Exec(query)
	if err == nil {
		// Дамп данных таблицы
		query = fmt.Sprintf(`INSERT INTO %s (id, login, password, name, userType, blocked) VALUES
			(1, 'admin', '1Jq0Kxw8KsnPNZNbf-QJZOed4BHo9uFeO', 'Никонов Виталий', 'admin', 0);
		`,
			table,
		)
		_, err = db.Exec(query)
	}
	return err
}

// FindManagerByChatID Поиск менеджера по ID чата
// func FindManagerByChatID(chatID int64, table string, db *sql.DB) (Manager, error) {
// 	s := Manager{ChatID: chatID}
// 	query := fmt.Sprintf(`SELECT id, userID, firstname, lastname, username, reghash, date, status FROM %s WHERE chatID = %d`, table, chatID)
// 	err := db.QueryRow(query).Scan(&s.ID, &s.UserID, &s.FirstName, &s.LastName, &s.UserName, &s.Reghash, &s.Datetime, &s.Status)
// 	return s, err
// }

// // Get Получение менеджера
// func Get(id int64, table string, db *sql.DB) (Manager, error) {
// 	s := Manager{ID: id}
// 	query := fmt.Sprintf(`SELECT userID, chatID, firstname, lastname, username, reghash, date, status FROM %s WHERE id = %d`, table, id)
// 	err := db.QueryRow(query).Scan(&s.UserID, &s.ChatID, &s.FirstName, &s.LastName, &s.UserName, &s.Datetime, &s.Status)
// 	return s, err
// }

// // NewUser Создание нового менеджера
// func NewUser(table string, db *sql.DB, chatID int64) (int64, error) {
// 	query := fmt.Sprintf(`INSERT INTO %s (chatID, reghash, date) VALUES (%d, '%s', DATETIMES)`, table, chatID, generate())
// 	query = strings.Replace(query, "DATETIMES", "strftime('%Y-%m-%d %H:%M:%S','now')", -1)
// 	res, err := db.Exec(query)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return res.LastInsertId()
// }

// // Update Новая комната
// func Update(user Manager, table string, db *sql.DB) error {
// 	_, err := db.Exec(fmt.Sprintf(`UPDATE %s SET userID=%d, firstname='%s', lastname='%s', username='%s', reghash='%s', status=%d WHERE chatID=%d`,
// 		table, user.UserID, user.FirstName, user.LastName, user.UserName, user.Reghash, user.Status, user.ChatID))
// 	return err
// }
// func generate() string {
// 	rand.Seed(time.Now().UnixNano())
// 	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
// 		"abcdefghijklmnopqrstuvwxyzåäö" +
// 		"0123456789")
// 	length := 6
// 	var b strings.Builder
// 	for i := 0; i < length; i++ {
// 		b.WriteRune(chars[rand.Intn(len(chars))])
// 	}
// 	str := b.String() // E.g. "ExcbsVQs"
// 	return str
// }
