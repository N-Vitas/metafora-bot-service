package restfull

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

// SettingsService Формирование сервиса настроек
func (app *Resource) SettingsService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/settings")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(app.GetSettings).
		Doc("Получение настроек").
		Operation("GetSettings"))
	ws.Route(ws.POST("/").To(app.UpdateSettings).
		Doc("Обновление настроек").
		Operation("UpdateSettings"))
	return ws
}

// GetSettings Эндпоинт всех настроек
func (app *Resource) GetSettings(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	r := Settings{ID: 1}
	query := `select token, updateID, comandID, crontime, googleFolder, hostService, durationManagers, 
		durationClients, durationStart, messageFailManager, messageFormAuth, fromMail, toMail, passMail, titleMail, bodyMail 
		from ` + app.Table("settings") + ` where id = 1`
	app.GetDb().QueryRow(query).Scan(
		&r.Token,
		&r.UpdateID,
		&r.ComandID,
		&r.Crontime,
		&r.GoogleFolder,
		&r.HostService,
		&r.DurationManagers,
		&r.DurationClients,
		&r.DurationStart,
		&r.MessageFailManager,
		&r.MessageFormAuth,
		&r.FromMail,
		&r.ToMail,
		&r.PassMail,
		&r.TitleMail,
		&r.BodyMail,
	)
	// Ответ пользователю
	WriteResponse(r, resp)
}

// UpdateSettings Обновление настроек
func (app *Resource) UpdateSettings(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	s := Settings{}
	decoder := json.NewDecoder(req.Request.Body)
	err := decoder.Decode(&s)
	if err != nil {
		WriteStatusError(http.StatusBadRequest, errors.New("Не удалось распарсить данные"), resp)
		return
	}
	if s.ID == 0 {
		WriteStatusError(http.StatusBadRequest, errors.New("Не верный ID настроек"), resp)
		return
	}
	_, err = app.GetDb().Exec(fmt.Sprintf(`UPDATE %s SET token='%s', updateID=%d, comandID=%d, crontime='%s', googleFolder='%s', 
	hostService='%s', durationManagers=%d, durationClients=%d, durationStart=%d, messageFailManager='%s', messageFormAuth='%s', fromMail='%s', toMail='%s', passMail='%s', titleMail='%s', bodyMail='%s' WHERE id=%d`,
		app.Table("settings"), s.Token, s.UpdateID, s.ComandID, s.Crontime, s.GoogleFolder, s.HostService,
		s.DurationManagers, s.DurationClients, s.DurationStart, s.MessageFailManager, s.MessageFormAuth, s.FromMail, s.ToMail, s.PassMail, s.TitleMail, s.BodyMail, s.ID))
	if err != nil {
		WriteStatusError(http.StatusInternalServerError, errors.New("Не удалось обновить настройки "+err.Error()), resp)
		return
	}
	WriteResponse(s, resp)
}

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
	FromMail		   string `json:"fromMail"`
	ToMail 			   string `json:"toMail"`
	PassMail 		   string `json:"passMail"`
	TitleMail 		   string `json:"titleMail"`
	BodyMail 		   string `json:"bodyMail"`
}
