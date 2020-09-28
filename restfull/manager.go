package restfull

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

// ManagerService Формирование сервиса групп
func (app *Resource) ManagerService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/manager")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(app.GetAllManager).
		Doc("Получение всех менеджеров").
		Operation("GetAllGroup"))
	ws.Route(ws.POST("/update").To(app.UpdateManager).
		Doc("Обновление менеджеров").
		Operation("UpdateManager"))
	ws.Route(ws.POST("/delete").To(app.DeleteManager).
		Doc("Удаление менеджеров").
		Operation("DeleteManager"))
	return ws
}

// GetAllManager Эндпоинт всех менеджеров
func (app *Resource) GetAllManager(req *restful.Request, resp *restful.Response) {
	managers := []Manager{}
	query := "select id, userID, chatID, firstname, lastname, username, reghash, status, date from " + app.Table("managers")
	rows, err := app.GetDb().Query(query)
	if err != nil {
		WriteStatusError(http.StatusForbidden, err, resp)
		return
	}
	for rows.Next() {
		ID := sql.NullInt64{}
		UserID := sql.NullInt64{}
		ChatID := sql.NullInt64{}
		FirstName := sql.NullString{}
		LastName := sql.NullString{}
		UserName := sql.NullString{}
		Reghash := sql.NullString{}
		Status := sql.NullInt64{}
		Datetime := sql.NullString{}
		err = rows.Scan(&ID, &UserID, &ChatID, &FirstName, &LastName, &UserName, &Reghash, &Status, &Datetime)
		if err != nil {
			continue
		}
		managers = append(managers, Manager{
			ID:        ID.Int64,
			UserID:    UserID.Int64,
			ChatID:    ChatID.Int64,
			FirstName: FirstName.String,
			LastName:  LastName.String,
			UserName:  UserName.String,
			Reghash:   Reghash.String,
			Status:    Status.Int64,
			Datetime:  Datetime.String,
		})
	}
	// Ответ пользователю
	WriteResponse(managers, resp)
}

// UpdateManager Обновление менеджеров
func (app *Resource) UpdateManager(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	manager := Manager{}
	decoder := json.NewDecoder(req.Request.Body)
	err := decoder.Decode(&manager)
	if err != nil {
		WriteStatusError(http.StatusBadRequest, errors.New("Не удалось распарсить данные"), resp)
		return
	}
	if manager.ID == 0 {
		WriteStatusError(http.StatusBadRequest, errors.New("Не верный ID менеджера"), resp)
		return
	}
	_, err = app.GetDb().Exec(fmt.Sprintf(`UPDATE %s SET firstname='%s', lastname='%s', username='%s', status=%d WHERE id=%d`,
		app.Table("managers"), manager.FirstName, manager.LastName, manager.UserName, manager.Status, manager.ID))
	if err != nil {
		WriteStatusError(http.StatusInternalServerError, errors.New("Не удалось обновить менеджера "+err.Error()), resp)
		return
	}
	query := "select userID, chatID, firstname, lastname, username, reghash, status, date from " + app.Table("managers") + " where id = ?"
	app.GetDb().QueryRow(query, manager.ID).Scan(&manager.UserID, &manager.ChatID, &manager.FirstName, &manager.LastName, &manager.UserName, &manager.Reghash, &manager.Status, &manager.Datetime)
	WriteResponse(manager, resp)
}

// DeleteManager Удаление менеджеров
func (app *Resource) DeleteManager(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	m := Manager{}
	decoder := json.NewDecoder(req.Request.Body)
	err := decoder.Decode(&m)
	if err != nil {
		WriteStatusError(http.StatusBadRequest, errors.New("Не удалось распарсить данные"), resp)
		return
	}
	if m.ID == 0 {
		WriteStatusError(http.StatusBadRequest, errors.New("Не верный ID менеджера"), resp)
		return
	}
	_, err = app.GetDb().Exec(fmt.Sprintf(`delete from %s WHERE id=%d`, app.Table("managers"), m.ID))
	if err != nil {
		WriteStatusError(http.StatusInternalServerError, errors.New("Не удалось удалить менеджера "+err.Error()), resp)
		return
	}
	_, err = app.GetDb().Exec(fmt.Sprintf(`delete from %s WHERE managerID=%d`, app.Table("group_manager"), m.ID))
	if err != nil {
		WriteStatusError(http.StatusInternalServerError, errors.New("Не удалось удалить менеджера "+err.Error()), resp)
		return
	}
	WriteSuccess(resp)
}

// Manager Структура Менеджера
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
