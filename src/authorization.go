package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var (
	SCOPE               = youtube.YoutubeReadonlyScope
	CLIENT_SECRETS_FILE string
	CREDENTIALS_FILE    string
)

func init() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatal("Error: A folder path must be given as a cli argument.")
	}
	CLIENT_SECRETS_FILE = filepath.Join(os.Args[1], "client_secrets.json")
	CREDENTIALS_FILE = filepath.Join(os.Args[1], "credentials.json")
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetToken(config *oauth2.Config) *oauth2.Token {
	tok, err := tokenFromFile(CREDENTIALS_FILE)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(CREDENTIALS_FILE, tok)
	}
	return tok
}

func GetConfig() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(CLIENT_SECRETS_FILE)
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secrets file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, SCOPE)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secrets file to config: %v", err)
	}
	return config, nil
}
