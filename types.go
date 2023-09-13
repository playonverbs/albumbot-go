package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"
)

const (
	dateLayout = "02/01/2006"
)

type ListenStatus int64

const (
	Listened ListenStatus = iota
	Listening
	NotListened
)

func (l ListenStatus) String() string {
	switch l {
	case Listened:
		return "Listened"
	case Listening:
		return "Currently Listening"
	case NotListened:
		return "Not Listened"

	default:
		return "Unknown"
	}
}

func NewListenStatus(status string) ListenStatus {
	switch status {
	case "Listened":
		return Listened
	case "Currently Listening":
		return Listening
	case "Not Listened":
		return NotListened

	default:
		return NotListened
	}
}

// Represents an Album score, rated out of 30
type Score uint8

func (s Score) String() string {
	return fmt.Sprintf("%d/30", s)
}

func parseScore(str string) (Score, error) {
	var score string
	_, err := fmt.Sscanf(str, "%d/30", score)
	if err != nil {
		return 0, nil
	}

	val, err := strconv.ParseUint(score, 10, 8)
	if err != nil {
		return 0, err
	}
	return Score(val), nil
}

type Entry interface {
	ToValue() []interface{}
}

// Entry contains fields to represent an Album row in the google spreadsheet
// TODO: consider just using string for SpotifyURL
type Album struct {
	ID          int // An ID only filled when reading from the sheet
	Album       string
	DateAdded   time.Time
	SuggestedBy string
	SpotifyURL  url.URL
	Votes       uint
	MeanScore   Score
	Status      ListenStatus
}

func (a Album) ToValue() []interface{} {
	dateFormat := a.DateAdded.Format(dateLayout)

	return []interface{}{
		a.Album,
		dateFormat,
		a.SuggestedBy,
		a.SpotifyURL.String(),
		a.Votes,
		a.MeanScore.String(),
		a.Status.String(),
	}
}

func NewAlbum(album, suggestedBy, spotifyURL string) *Album {
	u, err := url.Parse(spotifyURL)
	if err != nil {
		log.Printf("could not parse %s: %s", spotifyURL, err)
	}
	u.RawQuery = ""

	return &Album{
		Album:       album,
		DateAdded:   time.Now(),
		SuggestedBy: suggestedBy,
		SpotifyURL:  *u,
		Votes:       0,
		MeanScore:   0,
		Status:      NotListened,
	}
}

func NewAlbumFromRow(index int, row []interface{}) *Album {
	date, err := time.Parse(dateLayout, row[1].(string))
	if err != nil {
		log.Fatalf("cannot parse date: %#v", err)
	}

	u, err := url.Parse(row[3].(string))
	if err != nil {
		log.Fatalf("cannot parse url: %#v", err)
	}
	u.RawQuery = ""

	votes, err := strconv.ParseUint(row[4].(string), 10, 64)
	if err != nil {
		votes = 0
	}

	score, err := parseScore(row[5].(string))
	if err != nil {
		score = 0
	}

	return &Album{
		ID:          index,
		Album:       row[0].(string),
		DateAdded:   date,
		SuggestedBy: row[2].(string),
		SpotifyURL:  *u,
		Votes:       uint(votes),
		MeanScore:   score,
		Status:      NewListenStatus(row[6].(string)),
	}
}

type Albums []*Album

func (s Albums) Len() int      { return len(s) }
func (s Albums) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s Albums) AlbumsInWeek(date time.Time) Albums {
	var filtered Albums

	year, week := date.ISOWeek()
	for _, v := range s {
		thisYear, thisWeek := v.DateAdded.ISOWeek()
		if thisYear == year && thisWeek == week {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

type ByVotes struct{ Albums }

func (s ByVotes) Less(i, j int) bool { return s.Albums[i].Votes < s.Albums[j].Votes }

type ByDate struct{ Albums }

func (s ByDate) Less(i, j int) bool { return s.Albums[i].DateAdded.After(s.Albums[j].DateAdded) }

func CompareDates(d1, d2 time.Time) bool {
	ny, nm, nd := d1.Date()
	ey, em, ed := d2.Date()

	return ny == ey && nm == em && nd == ed
}

// func validateSpotifyLink(string URL)
