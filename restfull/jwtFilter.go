package restfull

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
)

// Auth Структура авторизации
type Auth struct {
	UserID    int64  `json:"id"`
	UserLogin string `json:"login"`
	UserName  string `json:"name"`
	UserType  string `json:"role"`
}

// JWTFilter Фильтр токена
func (s *Resource) JWTFilter(req *restful.Request) (Auth, error) {
	tokenHeader := req.Request.Header.Get("authorization")
	auth := Auth{}
	//Bearer
	if tokenHeader != "" {
		token, _ := jwt.Parse(strings.Replace(tokenHeader, "Bearer ", "", 1), func(token *jwt.Token) (interface{}, error) {

			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(s.secret), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			b, _ := json.Marshal(claims)
			json.Unmarshal(b, &auth)
			return auth, nil
		}
	}
	return auth, errors.New("Неверный токен авторизации")
}

// GetAuth Проверка токена
func (s *Resource) GetAuth(tokenHeader string) (Auth, error) {
	auth := Auth{}
	//Bearer
	if tokenHeader != "" {
		token, _ := jwt.Parse(strings.Replace(tokenHeader, "Bearer ", "", 1), func(token *jwt.Token) (interface{}, error) {

			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(s.secret), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			b, _ := json.Marshal(claims)
			json.Unmarshal(b, &auth)
			return auth, nil
		}
	}
	return auth, errors.New("Неверный токен авторизации")
}
