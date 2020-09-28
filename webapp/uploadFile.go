package webapp

import (
	"fmt"
	"metafora-bot-service/googledrive"
	"net/http"
)

// UploadFile Загрузка файла
func (s *ServerSoket) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Парсер данных формы, 10 << 20 максимальный вес файла 10 MB
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("clientFile")
	room := r.PostFormValue("chatRoom")
	if err != nil {
		Error("Error Retrieving the File %v", err)
		return
	}
	if link, ok := googledrive.CreateFile(s.srv, s.Controller.FolderID, handler, file); ok {
		Notice("CreateFile %s", link)
		if room, ok := s.Controller.GetRoom(room); ok {
			if room, ok = s.Controller.CreateClientMessage(link, room); ok {
				s.CheckClient(room.Room.Room, room.Room.ChatID, link, room)
			}
		}
		fmt.Fprintf(w, "Success Uploaded File\n")
		return
	}
	fmt.Fprintf(w, "Fail Uploaded File\n")
}
