package groups

import (
	"database/sql"
	"fmt"
	"strings"
)

// Group Структура группы
type Group struct {
	ID       int64  `json:"id"`
	ParentID int64  `json:"parentID"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Status   int64  `json:"status"`
	Datetime string `json:"date"`
	View     bool   `json:"view"`
}

// Init Создание базы данных для сообщений
func Init(table string, db *sql.DB) error {
	// query := fmt.Sprintf(`CREATE TABLE %s (
	// 	id int(9) NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	parentID int(20) NOT NULL DEFAULT 0,
	// 	name varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
	// 	title varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
	// 	view tinyint(1) NOT NULL DEFAULT 1,
	// 	date datetime NOT NULL DEFAULT current_timestamp(),
	// 	status int(5) NOT NULL DEFAULT 1
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	// 	table,
	// )
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		parentID integer NOT NULL DEFAULT 0,
		name text DEFAULT '',
		title text DEFAULT '',
		view integer NOT NULL DEFAULT 1,
		date text NOT NULL DEFAULT '',
		status integer NOT NULL DEFAULT 1
	  );`,
		table,
	)
	if _, err := db.Exec(query); err == nil {
		// Дамп данных таблицы
		query = fmt.Sprintf(`INSERT INTO %s (id, parentID, name, title, view, date, status) VALUES
		(1, 0, 'almaty', 'Алматы', 1, DATETIMES, 1),
		(2, 1, 'managers_one', 'Менеджеры 1', 0, DATETIMES, 1),
		(3, 1, 'managers_two', 'Менеджеры 2', 0, DATETIMES, 1),
		(4, 0, 'kustanay', 'Кустанай', 1, DATETIMES, 1);
		`,
			table,
		)
		query = strings.Replace(query, "DATETIMES", "strftime('%Y-%m-%d %H:%M:%S','now')", -1)
		_, err = db.Exec(query)
		return err
	}
	return nil
}

// Get Получение группы
func Get(id int64, table string, db *sql.DB) (Group, error) {
	s := Group{ID: id}
	query := fmt.Sprintf(`SELECT parentID, name, title, view, date, status FROM %s WHERE id = '%d'`, table, id)
	err := db.QueryRow(query).Scan(&s.ParentID, &s.Name, &s.Title, &s.Status, &s.Datetime, &s.View)
	return s, err
}

// GetParents Получение группы
func GetParents(parentID int64, table string, db *sql.DB) ([]Group, error) {
	res := []Group{}
	query := fmt.Sprintf(`SELECT id, parentID ,name ,title ,view ,date ,status FROM %s where parentID = %d`, table, parentID)
	rows, err := db.Query(query)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		s := Group{}
		err = rows.Scan(&s.ID, &s.ParentID, &s.Name, &s.Title, &s.Status, &s.Datetime, &s.View)
		if err != nil {
			continue
		}
		res = append(res, s)
	}
	return res, err
}

// GetAll Получение всех групп
func GetAll(table string, db *sql.DB) (map[int64]Group, error) {
	res := make(map[int64]Group)
	query := fmt.Sprintf(`SELECT id, parentID ,name ,title ,view ,date ,status FROM %s`, table)
	rows, err := db.Query(query)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		s := Group{}
		err = rows.Scan(&s.ID, &s.ParentID, &s.Name, &s.Title, &s.Status, &s.Datetime, &s.View)
		if err != nil {
			continue
		}
		res[s.ID] = s
	}
	return res, err
}
