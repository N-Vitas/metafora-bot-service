package restfull

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
)

// GroupService Формирование сервиса групп
func (app *Resource) GroupService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/group")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(app.GetAllGroup).
		Doc("Получение всех групп").
		Operation("GetAllGroup"))
	ws.Route(ws.GET("/{id}").To(app.GetGroupByID).
		Doc("Получение группы по ID").
		Operation("GetAllGroup").Param(ws.PathParameter("id", "ID").DataType("int")))
	ws.Route(ws.POST("/create").To(app.CreateGroup).
		Doc("Создание группы").
		Operation("GetAllGroup"))
	ws.Route(ws.POST("/update").To(app.UpdateGroup).
		Doc("Обновление группы").
		Operation("UpdateGroup"))
	ws.Route(ws.POST("/delete").To(app.DeleteGroup).
		Doc("Удаление группы").
		Operation("DeleteGroup"))
	return ws
}

// GetAllGroup Эндпоинт всех групп
func (app *Resource) GetAllGroup(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	groups := []Group{}
	query := "select id, parentID, name, title, view, date, status from " + app.Table("groups")
	rows, err := app.GetDb().Query(query)
	if err != nil {
		WriteStatusError(http.StatusForbidden, err, resp)
		return
	}
	for rows.Next() {
		ID := sql.NullInt64{}
		ParentID := sql.NullInt64{}
		Name := sql.NullString{}
		Title := sql.NullString{}
		View := sql.NullInt64{}
		Date := sql.NullString{}
		Status := sql.NullInt64{}
		err = rows.Scan(&ID, &ParentID, &Name, &Title, &View, &Date, &Status)
		if err != nil {
			continue
		}
		managers := []int64{}
		query = "select managerID from " + app.Table("group_manager") + " where groupID = ?"
		makes, _ := app.GetDb().Query(query, ID.Int64)
		if err == nil {
			for makes.Next() {
				managerID := sql.NullInt64{}
				makes.Scan(&managerID)
				if managerID.Valid {
					managers = append(managers, managerID.Int64)
				}
			}
		}
		groups = append(groups, Group{
			ID:       ID.Int64,
			ParentID: ParentID.Int64,
			Name:     Name.String,
			Title:    Title.String,
			View:     View.Int64,
			Date:     Date.String,
			Status:   Status.Int64,
			Managers: managers,
		})
	}
	// Ответ пользователю
	WriteResponse(groups, resp)
}

// GetGroupByID Эндпоинт группы по ID
func (app *Resource) GetGroupByID(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	ID, err := strconv.Atoi(req.PathParameter("id"))
	group := Group{}
	if err != nil {
		WriteStatusError(http.StatusBadRequest, errors.New("Нет паратеметра id"), resp)
		return
	}
	query := "select id, parentID, name, title, view, date, status from " + app.Table("groups") +
		" where id = ?"
	err = app.GetDb().QueryRow(query, ID).Scan(&group.ID, &group.ParentID, &group.Name, &group.Title, &group.View, &group.Date, &group.Status)
	if err != nil {
		WriteStatusError(http.StatusNotFound, errors.New("Нет паратеметра id"), resp)
		return
	}
	managers := []int64{}
	query = "select managerID from " + app.Table("group_manager") + " where groupID = ?"
	makes, _ := app.GetDb().Query(query, group.ID)
	if err == nil {
		for makes.Next() {
			managerID := sql.NullInt64{}
			makes.Scan(&managerID)
			if managerID.Valid {
				managers = append(managers, managerID.Int64)
			}
		}
	}
	group.Managers = managers
	WriteResponse(group, resp)
}

// CreateGroup Создание новой группы
func (app *Resource) CreateGroup(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	group := Group{}
	decoder := json.NewDecoder(req.Request.Body)
	err := decoder.Decode(&group)
	if err != nil {
		Info("%v", err)
		WriteStatusError(http.StatusBadRequest, errors.New("Не удалось распарсить данные"), resp)
		return
	}
	if len(group.Title) == 0 {
		WriteStatusError(http.StatusBadRequest, errors.New("Поле title не должно быть пустым"), resp)
		return
	}
	id := sql.NullInt64{}
	err = app.GetDb().QueryRow("select max(id) from " + app.Table("groups")).Scan(&id)
	if err != nil {
		WriteStatusError(http.StatusNotFound, errors.New("Нет паратеметра id"), resp)
		return
	}
	group.ID = id.Int64 + 1
	group.Name = Transcript(group.Title)
	query := fmt.Sprintf(`INSERT INTO %s (id, name, title, view, date) VALUES (%d, '%s', '%s', %d, DATETIMES)`,
		app.Table("groups"), group.ID, group.Name, group.Title, group.View)
	query = strings.Replace(query, "DATETIMES", "strftime('%Y-%m-%d %H:%M:%S','now')", -1)
	_, err = app.GetDb().Exec(query)
	if err != nil {
		WriteStatusError(http.StatusInternalServerError, errors.New("Не удалось создать группу "+err.Error()), resp)
		return
	}
	app.FixGroupManagers(group)
	query = "select parentID, name, title, view, date, status from " + app.Table("groups") + " where id = ?"
	app.GetDb().QueryRow(query, group.ID).Scan(&group.ParentID, &group.Name, &group.Title, &group.View, &group.Date, &group.Status)
	app.FixBotTell()
	WriteResponse(group, resp)
}

// UpdateGroup Обновление группы
func (app *Resource) UpdateGroup(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	group := Group{}
	decoder := json.NewDecoder(req.Request.Body)
	err := decoder.Decode(&group)
	if err != nil {
		WriteStatusError(http.StatusBadRequest, errors.New("Не удалось распарсить данные"), resp)
		return
	}
	if group.ID == 0 {
		WriteStatusError(http.StatusBadRequest, errors.New("Не верный ID группы"), resp)
		return
	}
	if len(group.Title) == 0 {
		WriteStatusError(http.StatusBadRequest, errors.New("Поле title не должно быть пустым"), resp)
		return
	}
	group.Name = Transcript(group.Title)
	_, err = app.GetDb().Exec(fmt.Sprintf(`UPDATE %s SET parentID=%d, name='%s', title='%s', view=%d, status=%d WHERE id=%d`,
		app.Table("groups"), group.ParentID, group.Name, group.Title, group.View, group.Status, group.ID))
	if err != nil {
		WriteStatusError(http.StatusInternalServerError, errors.New("Не удалось обновить группу "+err.Error()), resp)
		return
	}
	app.FixGroupManagers(group)
	query := "select parentID, name, title, view, date, status from " + app.Table("groups") + " where id = ?"
	app.GetDb().QueryRow(query, group.ID).Scan(&group.ParentID, &group.Name, &group.Title, &group.View, &group.Date, &group.Status)
	app.FixBotTell()
	WriteResponse(group, resp)
}

// DeleteGroup Удаление группы
func (app *Resource) DeleteGroup(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	group := Group{}
	decoder := json.NewDecoder(req.Request.Body)
	err := decoder.Decode(&group)
	if err != nil {
		WriteStatusError(http.StatusBadRequest, errors.New("Не удалось распарсить данные"), resp)
		return
	}
	if group.ID == 0 {
		WriteStatusError(http.StatusBadRequest, errors.New("Не верный ID группы"), resp)
		return
	}
	_, err = app.GetDb().Exec(fmt.Sprintf(`delete from %s WHERE id=%d`, app.Table("groups"), group.ID))
	if err != nil {
		WriteStatusError(http.StatusInternalServerError, errors.New("Не удалось удалить группу "+err.Error()), resp)
		return
	}
	app.FixBotTell()
	WriteSuccess(resp)
}

// FixGroupManagers Синхронизация менеджеров в группе
func (app *Resource) FixGroupManagers(group Group) {
	app.GetDb().Exec("delete from "+app.Table("group_manager")+" where groupID=?", group.ID)
	if len(group.Managers) > 0 {
		for _, id := range group.Managers {
			_, err := app.GetDb().Exec("INSERT INTO "+app.Table("group_manager")+" (groupID, managerID) VALUES (?, ?)", group.ID, id)
			if err != nil {
				Info("%v", err)
			}
		}
	}
}

// FixBotTell Синхронизация в реплики бота
func (app *Resource) FixBotTell() {
	query := "select title from " + app.Table("groups") + " where view = 1"
	rows, err := app.GetDb().Query(query)
	if err != nil {
		Info("%v", err)
		return
	}
	msg := []string{}
	for rows.Next() {
		res := sql.NullString{}
		rows.Scan(&res)
		if res.Valid {
			msg = append(msg, res.String)
		}
	}
	_, err = app.GetDb().Exec(fmt.Sprintf(`UPDATE %s SET dataType='%s' WHERE type = 'select'`, app.Table("bot"), strings.Join(msg, ",")))
	if err != nil {
		Info("%v", err)
		return
	}
}

// Group Структура группы
type Group struct {
	ID       int64   `json:"id"`
	ParentID int64   `json:"parentID"`
	Name     string  `json:"name"`
	Title    string  `json:"title"`
	View     int64   `json:"view"`
	Date     string  `json:"date"`
	Status   int64   `json:"status"`
	Managers []int64 `json:"managers"`
}
