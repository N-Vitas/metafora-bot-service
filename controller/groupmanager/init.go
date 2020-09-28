package groupmanager

import (
	"database/sql"
	"fmt"
	"metafora-bot-service/controller/managers"
)

// GroupManager Структура Связи менеджера и группы
type GroupManager struct {
	GroupID   int64 `json:"groupID"`
	ManagerID int64 `json:"managerID"`
}

// Init Создание базы данных для сообщений
func Init(table string, db *sql.DB, group string, manager string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		groupID int(20) DEFAULT 0,
		managerID int(20) DEFAULT 0,
		CONSTRAINT group_id_ibfk_1 FOREIGN KEY (groupID) REFERENCES %s (id),
		CONSTRAINT manager_id_ibfk_2 FOREIGN KEY (managerID) REFERENCES %s (id)
	  );
	  `,
		table,
		group,
		manager,
	)
	_, err := db.Exec(query)
	return err
}

// GetAll Получение всех менеджеров по группе
func GetAll(groupID int64, manager, group string, db *sql.DB) (map[int64]managers.Manager, error) {
	res := make(map[int64]managers.Manager)
	query := fmt.Sprintf(`SELECT m.id, m.userID, m.chatID, m.firstname, m.lastname, m.username, m.reghash, m.date, m.status FROM %s m 
		INNER JOIN %s gm ON gm.managerID = m.Id WHERE gm.groupID = %d`, manager, group, groupID)
	rows, err := db.Query(query)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		s := managers.Manager{}
		err = rows.Scan(&s.ID, &s.UserID, &s.ChatID, &s.FirstName, &s.LastName, &s.UserName, &s.Reghash, &s.Datetime, &s.Status)
		if err != nil {
			continue
		}
		res[s.ID] = s
	}
	return res, err
}
