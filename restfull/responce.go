package restfull

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

// Error Структура ответа ошибки
type Error struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

// Response Структура ответа
type Response struct {
	Success bool        `json:"success"`
	Error   *Error      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteSuccess Пустой успешный ответ
func WriteSuccess(resp *restful.Response) {
	NewResponse(true).WriteStatus(200, resp)
}

// WriteResponse Успешный ответ с данными
func WriteResponse(data interface{}, resp *restful.Response) {
	WriteResponseStatus(200, data, resp)
}

// WriteResponseStatus Успешный ответ с данными и кодом ответа
func WriteResponseStatus(status int, data interface{}, resp *restful.Response) {
	success := NewResponse(true)
	success.Data = data
	success.WriteStatus(status, resp)
}

// NewResponse Формирование структуры успешного ответа
func NewResponse(success bool) *Response {
	return &Response{Success: success}
}

// NewErrorResponse Формирование структуры ответа ошибки
func NewErrorResponse(err error) *Response {
	res := &Response{Success: false, Error: &Error{}}
	res.SetError(err)
	return res
}

// SetError Формирование структуры ошибки
func (r *Response) SetError(err error) {
	if err != nil {
		if r.Error == nil {
			r.Error = &Error{}
		}
		r.Error.Name = err.Error()
	}
}

// WriteStatus Ответ кода
func (r *Response) WriteStatus(status int, resp *restful.Response) {
	if r.Error != nil && r.Error.Code == 0 {
		r.Error.Code = status
	}
	resp.WriteHeaderAndEntity(status, r)
}

// WriteError Ответ кода ошибки
func WriteError(err error, resp *restful.Response) {

	// Set response status code
	code := http.StatusInternalServerError

	// String error
	error := err.Error()

	if error == "not found" || len(error) > 7 && error[:7] == "Unknown" {
		code = http.StatusNotFound
	} else if error == "unauthorized" || len(error) > 14 && error[:14] == "not authorized" {
		code = http.StatusUnauthorized
	}

	// Write error response
	WriteStatusError(code, err, resp)
}

// WriteStatusError Ответ кода ошибки с данными ошибки
func WriteStatusError(status int, err error, resp *restful.Response) {
	app := NewResponse(false)
	app.SetError(err)
	app.WriteStatus(status, resp)
}
