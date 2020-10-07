package restfull

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

// StatisticService Формирование сервиса статистики
func (app *Resource) StatisticService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/statistic")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(app.GetAllStats).
		Doc("Получение статистики").
		Operation("GetAllStats"))
	return ws
}

// GetAllStats Эндпоинт всей статистики
func (app *Resource) GetAllStats(req *restful.Request, resp *restful.Response) {
	_, forbiden := app.JWTFilter(req)
	if forbiden != nil {
		WriteStatusError(http.StatusUnauthorized, forbiden, resp)
		return
	}
	r := Statistic{}
	query := `select * from (
		(select count(r1.id) openChats from ` + app.Table("rooms") + ` r1 where r1.status > 0),
		(select count(r2.id) closedChats from ` + app.Table("rooms") + ` r2 where r2.status = 0),
		(select count(m1.id) managers from ` + app.Table("managers") + ` m1 where m1.status > 1),
		(select count(m2.id) guests from ` + app.Table("managers") + ` m2 where m2.status = 1)
		)`
	app.GetDb().QueryRow(query).Scan(&r.OpenChats, &r.CloseChats, &r.Managers, &r.Guests)
	// Ответ пользователю
	WriteResponse(r, resp)
}

// Statistic Структура статистики
type Statistic struct {
	Managers   int64 `json:"managers"`
	Guests     int64 `json:"guests"`
	OpenChats  int64 `json:"openChats"`
	CloseChats int64 `json:"closeChats"`
}
