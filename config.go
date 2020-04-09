package go_gdrive

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	GOOGLE_DIRECTORY             = "application/vnd.google-apps.folder"
	GOOGLE_SPREAD_SHEET          = "application/vnd.google-apps.spreadsheet"
	GOOGLE_DOCS                  = "application/vnd.google-apps.document"
	GOOGLE_DIRECTORY_Q_MIME_TYPE = "mimeType='application/vnd.google-apps.folder' and trashed=false"
)

type googleApiInterface interface {
	GetDriveServices() (gdriveInterface, error)
	GetSpreadSheetServices() (SpreadSheetInterface, error)
	RefreshToken()(NewToken *oauth2.Token,err error)
}
type googleApi struct {
	client    *http.Client
	config    *oauth2.Config
	tokenPath string
	userToken *oauth2.Token
}

func Conenction(CredentialsPath string, TokenPath string) (googleApiInterface, error) {
	credentials, err := ioutil.ReadFile(CredentialsPath)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to read credentials json file. Err: %v\n", err))
		return nil, err
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credentials, drive.DriveFileScope)

	// this function will run for the first time when you dont have token file
	Token, err := tokenFromFile(TokenPath)
	if err != nil {
		Token = getTokenFromWeb(config)
		saveToken(TokenPath, Token)
	}
	gdrive := googleApi{
		client: config.Client(context.Background(), Token),
		config:config,
		tokenPath:TokenPath,
		userToken:Token,
	}
	return &gdrive, err
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

func (g googleApi)RefreshToken()(NewToken *oauth2.Token,err error){
	tokenSource := g.config.TokenSource(oauth2.NoContext, g.userToken)
	NewToken, err = tokenSource.Token()
	if err != nil {
		return nil, err
	}
	saveToken(g.tokenPath,NewToken)
	return NewToken,err
}