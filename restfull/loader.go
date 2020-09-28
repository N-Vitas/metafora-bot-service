package restfull

import "github.com/emicklei/go-restful"

// Loader Загружает список сервисов
func (app *Resource) Loader() {
	app.Add(app.LoginService())
	app.Add(app.GroupService())
	app.Add(app.ManagerService())
}

// Add добавляет в список сервис
func (app *Resource) Add(service *restful.WebService) {
	app.routes = append(app.routes, service)
}
