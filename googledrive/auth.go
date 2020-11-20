package googledrive

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	// "path/filepath"

	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Auth подключение гугл диска
func Auth() *drive.Service {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return srv
}

// ReadGoogleDrive 123
func ReadGoogleDrive(srv *drive.Service) {
	// srv.Changes.GetStartPageToken().DriveId("1Jq0Kxw8KsnPNZNbf-QJZOed4BHo9uFeO")
	// = "1Jq0Kxw8KsnPNZNbf-QJZOed4BHo9uFeO"

	r, err := srv.Files.List().PageSize(10).SupportsTeamDrives(true).
		Fields("nextPageToken, files(id, name)").Q("'1MmY3uvmNcOk1Zg_9gaW2KMd7qly2H2WQ' in parents").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
}

// ClearGoogleDrive 123
func ClearGoogleDrive(srv *drive.Service) {
	// srv.Changes.GetStartPageToken().DriveId("1Jq0Kxw8KsnPNZNbf-QJZOed4BHo9uFeO")
	// = "1Jq0Kxw8KsnPNZNbf-QJZOed4BHo9uFeO"

	r, err := srv.Files.List().PageSize(10).SupportsTeamDrives(true).
		Fields("nextPageToken, files(id, name)").Q("'1MmY3uvmNcOk1Zg_9gaW2KMd7qly2H2WQ' in parents").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
			srv.Files.Delete(i.Id).Do()

		}
	}
}

// CreateFile Создание файла в гугл диске
func CreateFile(/*srv *drive.Service, folder string, */handler *multipart.FileHeader, file multipart.File) (string, bool) {
	f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	defer func() {
		file.Close()
		f.Close()
	}()
	io.Copy(f, file)
	f.Close()
	return fmt.Sprintf("http://185.144.29.230:8082/public/%s", handler.Filename), true
	// f, err = os.Open("./upload/" + handler.Filename)
	// fileInfo, err := f.Stat()
	// if err != nil {
	// 	fmt.Println("Stat", err)
	// 	return "", false
	// }
	// inFile := &drive.File{
	// 	Name:     filepath.Base(fileInfo.Name()),
	// 	Parents:  []string{folder},
	// 	MimeType: "application/octet-stream",
	// }
	// outFile, err := srv.Files.Create(inFile).Media(f).Do()
	// if err != nil {
	// 	fmt.Println("CreateFile err", err)
	// 	return "", false
	// }
	// return fmt.Sprintf("https://drive.google.com/file/d/%s/view", outFile.Id), true
}
