package go_gdrive

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)


const (
	GOOGLE_DIRECTORY  = "application/vnd.google-apps.folder"
	GOOGLE_SPREAD_SHEET = "application/vnd.google-apps.spreadsheet"
	GOOGLE_DOCS = "application/vnd.google-apps.document"
	GOOGLE_DIRECTORY_Q_MIME_TYPE = "mimeType='application/vnd.google-apps.folder' and trashed=false"
)

type googleApiInterface interface {
	GetDriveServices()(gdriveInterface,error)
	GetSpreadSheetServices()(SpreadSheetInterface,error)
}
type googleApi struct {
	client *http.Client
}

func Conenction(CredentialsPath string)(googleApiInterface, error){
	credentials, err := ioutil.ReadFile(CredentialsPath)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to read credentials json file. Err: %v\n", err))
		return nil, err
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credentials, drive.DriveFileScope)

	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	gdrive := googleApi{
		client:config.Client(context.Background(), tok),
	}
	return &gdrive,err
}

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
func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
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

