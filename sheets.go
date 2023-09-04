package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Service struct{ *sheets.Service }

func NewService(ctx context.Context, opt option.ClientOption) (*Service, error) {
	srv, err := sheets.NewService(ctx, opt)
	if err != nil {
		return nil, err
	}
	return &Service{srv}, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
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
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
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

// Get a list of entries from the google sheet and given range
// func GetSheetEntries(srv *sheets.Service, sheetID string, readRange string) ([]*Entry, error) {
func (srv Service) GetSheetEntries(sheetID string, readRange string) ([]*Entry, error) {
	var entries []*Entry

	res, err := srv.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		log.Printf("could not read sheet rows: %v", err)
		return nil, err
	}

	if len(res.Values) == 0 {
		log.Printf("No data found")
		return []*Entry{}, nil
	} else {
		for idx, row := range res.Values {
			entries = append(entries, NewEntryFromRow(idx, row))
		}
	}

	return entries, nil
}

// Append a single entry to the google sheet
func (srv Service) AppendSheetEntry(sheetID string, writeRange string, entry *Entry) error {
	body := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          writeRange,
	}
	body.Values = append(body.Values, entry.ToValue())

	res, err := srv.Spreadsheets.Values.Append(sheetID, writeRange, body).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Do()
	if err != nil {
		log.Printf("Unable to write data to sheet: %v", err)
		return err
	}

	log.Printf("Entry written with response %d", res.HTTPStatusCode)
	return nil
}
