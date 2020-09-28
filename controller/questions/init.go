package questions

import (
	"database/sql"
	"fmt"
)

// Replic Структура вопросов бота
type Replic struct {
	ID       int64  `json:"id"`
	Message  string `json:"message"`
	Type     string `json:"type"`
	DataType string `json:"dataType"`
	Sort     int64  `json:"sort"`
	Date     string `json:"date"`
	Status   int64  `json:"status"`
}

// Init Создание базы данных для сообщений
func Init(table string, db *sql.DB) error {
	// query := fmt.Sprintf(`CREATE TABLE %s (
	// 	id int(9) NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	message varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
	// 	type varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT 'text',
	// 	dataType varchar(100) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '',
	// 	sort int(20) NOT NULL DEFAULT 0,
	// 	date datetime NOT NULL DEFAULT current_timestamp(),
	// 	status int(5) NOT NULL DEFAULT 1
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	// 	table,
	// )
	query := fmt.Sprintf(`CREATE TABLE %s (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		message TEXT NOT NULL,
		type TEXT DEFAULT 'text',
		dataType TEXT NOT NULL DEFAULT '',
		sort integer NOT NULL DEFAULT 0,
		date TEXT NOT NULL DEFAULT '',
		status integer NOT NULL DEFAULT 1
	  );`,
		table,
	)
	if _, err := db.Exec(query); err == nil {
		// Дамп данных таблицы
		query = `INSERT INTO ` + table + ` (id, message, type, dataType, sort, date, status) VALUES
			(1, 'Добый день. Я могу Вам чем нибудь помочь?', 'text', '', 0, strftime('%Y-%m-%d %H:%M:%S','now'), 1),
			(2, 'Выберите Ваш город', 'select', 'Алматы,Кустанай', 1, strftime('%Y-%m-%d %H:%M:%S','now'), 1),
			(3, 'Что желаете перевести?', 'button', 'Текст,Личные документы', 2, strftime('%Y-%m-%d %H:%M:%S','now'), 1),
			(4, 'Заверить перевод у нотариуса?', 'button', 'Да,Нет', 3, strftime('%Y-%m-%d %H:%M:%S','now'), 1),
			(5, 'Вы желаете получить скидку или Вам нужно срочно?', 'button', 'Хочу скидку, Срочно', 4, strftime('%Y-%m-%d %H:%M:%S','now'), 1),
			(6, 'Можете прикрепить файл?', 'file', 'Прикрепить файл, Нет', 5, strftime('%Y-%m-%d %H:%M:%S','now'), 1);
		`
		_, err = db.Exec(query)
		return err
	}
	return nil
}

// Get Получение реплики
func Get(id int64, table string, db *sql.DB) (Replic, error) {
	s := Replic{ID: id}
	query := fmt.Sprintf(`SELECT message, type, dataType, sort, date, status FROM %s WHERE id = '%d'`, table, id)
	err := db.QueryRow(query).Scan(&s.Message, &s.Type, &s.DataType, &s.Sort, &s.Date, &s.Status)
	return s, err
}

// GetAll Получение всех реплик
func GetAll(table string, db *sql.DB) []Replic {
	s := []Replic{}
	query := fmt.Sprintf(`SELECT id, message, type, dataType, sort, date, status FROM %s ORDER BY sort ASC`, table)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("GetAll Replic error ", err)
		return s
	}
	defer rows.Close()
	for rows.Next() {
		r := Replic{}
		err = rows.Scan(&r.ID, &r.Message, &r.Type, &r.DataType, &r.Sort, &r.Date, &r.Status)
		if err != nil {
			continue
		}
		s = append(s, r)
	}
	return s
}

// Next Получение реплики
func Next(sort int64, table string, db *sql.DB) (Replic, error) {
	s := Replic{Sort: sort + 1}
	query := fmt.Sprintf(`SELECT id, message, type, dataType, sort, date, status FROM %s WHERE sort = %d`, table, sort+1)
	err := db.QueryRow(query).Scan(&s.ID, &s.Message, &s.Type, &s.DataType, &s.Sort, &s.Date, &s.Status)
	return s, err
}

// MaxSort Получение последней сортировки реплики
func MaxSort(table string, db *sql.DB) int64 {
	res := sql.NullInt64{}
	query := fmt.Sprintf(`SELECT max(sort) FROM %s`, table)
	db.QueryRow(query).Scan(&res)
	return res.Int64
}
