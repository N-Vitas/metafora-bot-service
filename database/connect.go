package database

import (
	"database/sql"
	"sync"

	// _ "github.com/go-sql-driver/mysql" // драйвер базы данных
	_ "github.com/mattn/go-sqlite3" // драйвер базы данных
)

// SessionDb Структура базы данных
type SessionDb struct {
	db         *sql.DB
	accessLock *sync.RWMutex
}

// New Инициализация сессии БД.
func New() *SessionDb {
	return &SessionDb{
		accessLock: &sync.RWMutex{},
	}
}

// GetDb Получение подключения к базе
func (s *SessionDb) GetDb() *sql.DB {
	s.accessLock.RLock()
	existing := s.db
	s.accessLock.RUnlock()
	// Возврат существующего подключения
	if existing != nil {
		err := existing.Ping()
		if err != nil {
			existing = nil
		} else {
			return existing
		}
	}
	s.accessLock.Lock()
	var newSession *sql.DB
	// conf, err := config.New()
	// if err != nil {
	// 	panic("Cannot read config file" + err.Error())
	// }
	// connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, password, server, port, config.GetString("database", "invo"))
	// connString := fmt.Sprintf("root:123@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	// connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
	// 	conf.DBUSER, conf.DBPASS, conf.DBHOST, conf.DBPORT, conf.DBNAME)
	// newSession, err = sql.Open("mysql", connString)
	newSession, err := sql.Open("sqlite3", "sqlite-database.db")
	if err != nil {
		panic("Cannot connect to database" + err.Error())
	}
	s.db = newSession
	s.accessLock.Unlock()
	return newSession
}
