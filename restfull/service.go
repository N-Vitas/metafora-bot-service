package restfull

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
)

// Resource Структура пакета
type Resource struct {
	connections func() *sql.DB
	tables      func(name string) string
	routes      []*restful.WebService
	secret      string
}

// Init инициализация апи
func Init(db func() *sql.DB, tables func(name string) string, secret string, cors bool) *Resource {
	app := &Resource{
		connections: db,
		routes:      []*restful.WebService{},
		secret:      secret,
		tables:      tables,
	}
	app.Loader()
	app.Register(cors)
	return app
}

// GetDb Получение подключения к базе
func (s *Resource) GetDb() *sql.DB {
	return s.connections()
}

// Table Получение названия таблицы
func (s *Resource) Table(name string) string {
	return s.tables(name)
}

// Register Регистрация сервисов
func (s *Resource) Register(cors bool) {
	app := restful.DefaultContainer
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	// gzip if accepted
	app.EnableContentEncoding(true)
	// faster router
	app.Router(restful.CurlyRouter{})
	if cors {
		corsRule := restful.CrossOriginResourceSharing{
			//ExposeHeaders: []string{"Content-Type"},
			AllowedDomains: []string{"http://127.0.0.1", "http://localhost"},
			AllowedHeaders: []string{"*", "content-type", "authorization", "uuid", "Accept", "X-Custom-Header", "Origin", "Access-Control-Allow-Origin"},
			AllowedMethods: []string{"POST", "GET", "OPTIONS"},
			CookiesAllowed: false,
			Container:      app,
		}
		app.Filter(corsRule.Filter)
	}
	for _, route := range s.routes {
		app.Add(route)
	}
}

// Info Вывод информации в консоль
func Info(template string, values ...interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf(t+" \033[1;30m[restfull][info]\033[0m \033[1;33m"+template+"\033[0m\n", values...)
}

// Transcript Транскрипт кирилицы в латиницу
func Transcript(str string) string {
	transcript := func(r rune) rune {
		switch {
		case r >= 'А' && r <= 'Я':
			return 'A' + (r-'A'+13)%26
		case r >= 'а' && r <= 'я':
			return 'a' + (r-'a'+13)%26
		case r == ' ':
			return '_'
		}
		return r
	}
	return strings.Map(transcript, strings.ToLower(str))
}
