package restfull

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
)

// LoginService Формирование сервиса авторизации
func (app *Resource) LoginService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/login")
	ws.Consumes("*/*")
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.POST("/token").To(app.GetToken).
		Doc("Get token by login credentials").
		Operation("GetToken"))
	return ws
}

// GetToken Эндпоинт авторизации
func (app *Resource) GetToken(req *restful.Request, resp *restful.Response) {
	var (
		UserID      sql.NullInt64
		UserLogin   sql.NullString
		UserName    sql.NullString
		UserType    sql.NullString
		UserBlocked sql.NullInt64
	)
	auth := AuthRequest{}
	decoder := json.NewDecoder(req.Request.Body)
	err := decoder.Decode(&auth)
	if err != nil {
		WriteStatusError(http.StatusBadRequest, errors.New("Не удалось распарсить данные"), resp)
		return
	}
	Info("%s", GetMD5Hash(auth.Password))
	query := "select id, login, name, userType, blocked from " + app.Table("user") + " where login = ? and password = ?"
	Info("%s", GetMD5Hash(auth.Password))
	err = app.GetDb().QueryRow(query, auth.Login, GetMD5Hash(auth.Password)).Scan(&UserID, &UserLogin, &UserName, &UserType, &UserBlocked)
	if err != nil {
		WriteStatusError(http.StatusForbidden, err, resp)
		// WriteStatusError(http.StatusForbidden, errors.New("Неверный логин или пароль"), resp)
		return
	}
	if UserID.Valid == false {
		WriteStatusError(http.StatusForbidden, errors.New("Неверный логин или пароль"), resp)
		return
	}
	// Проверка блокировки
	if UserBlocked.Int64 > 0 {
		WriteStatusError(http.StatusForbidden, errors.New("Ваша учетная запись заблокирована"), resp)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    UserID.Int64,
		"login": UserLogin.String,
		"name":  UserName.String,
		"role":  UserType.String,
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(app.secret))
	if err != nil {
		WriteStatusError(http.StatusForbidden, errors.New("Неверный логин или пароль"), resp)
		return
	}
	// Write response back to client
	WriteResponse(map[string]interface{}{
		"token": tokenString,
		"id":    UserID.Int64,
		"login": UserLogin.String,
		"name":  UserName.String,
		"role":  UserType.String,
	}, resp)
}

// GetMD5Hash Хеширование строки
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// AuthRequest Структура запроса авторизации
type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
