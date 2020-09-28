package settings

import (
	"database/sql"
	"fmt"
)

// Settings Структура настроек
type Settings struct {
	ID                 int64  `json:"id"`
	Token              string `json:"token"`
	UpdateID           int64  `json:"updateID"`
	ComandID           int64  `json:"comandID"`
	Crontime           string `json:"crontime"`
	GoogleFolder       string `json:"googleFolder"`
	HostService        string `json:"hostService"`
	DurationManagers   int64  `json:"durationManagers"`
	DurationClients    int64  `json:"durationClients"`
	DurationStart      int64  `json:"durationStart"`
	MessageFailManager string `json:"messageFailManager"`
	MessageFormAuth    string `json:"messageFormAuth"`
}

// Init Создание базы данных для сообщений
func Init(table string, db *sql.DB) error {
	// query := fmt.Sprintf(`CREATE TABLE %s (
	// 	id int(9) NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	token varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
	// 	updateID int(10) NOT NULL DEFAULT 0,
	// 	comandID int(10) NOT NULL DEFAULT 0,
	// 	crontime varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT 'hourly',
	// 	googleFolder varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	// 	hostService varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	// 	durationManagers int(10) NOT NULL DEFAULT 5,
	// 	durationClients int(10) NOT NULL DEFAULT 5,
	// 	durationStart int(10) NOT NULL DEFAULT 30,
	// 	messageFailManager varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	// 	messageFormAuth varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	// 	table,
	// )
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		token TEXT,
		updateID integer NOT NULL DEFAULT 0,
		comandID integer NOT NULL DEFAULT 0,
		crontime TEXT,
		googleFolder TEXT,
		hostService TEXT,
		durationManagers integer NOT NULL DEFAULT 5,
		durationClients integer NOT NULL DEFAULT 5,
		durationStart integer NOT NULL DEFAULT 30,
		messageFailManager TEXT,
		messageFormAuth TEXT
	  );`,
		table,
	)
	if _, err := db.Exec(query); err == nil {
		// Дамп данных таблицы
		query = fmt.Sprintf(`INSERT INTO %s (id, token, crontime, updateID, comandID, googleFolder, hostService, durationManagers, durationClients, durationStart, messageFailManager, messageFormAuth) VALUES
			(1, 'token', '10_sec', 0, 0, '1Jq0Kxw8KsnPNZNbf-QJZOed4BHo9uFeO', '0.0.0.0:8082', 30, 30, 300, 'messageFailManager','messageFormAuth');
		`,
			table,
		)
		_, err = db.Exec(query)
		return err
	}
	return nil
}

// Get Получение настроек
func Get(table string, db *sql.DB) (Settings, error) {
	s := Settings{}
	query := fmt.Sprintf("SELECT id, token, updateID, comandID, crontime, googleFolder, hostService, durationManagers, durationClients, durationStart, messageFailManager, messageFormAuth FROM %s WHERE id = 1", table)
	err := db.QueryRow(query).Scan(&s.ID, &s.Token, &s.UpdateID, &s.ComandID, &s.Crontime, &s.GoogleFolder, &s.HostService, &s.DurationManagers, &s.DurationClients, &s.DurationStart, &s.MessageFailManager, &s.MessageFormAuth)
	return s, err
}
