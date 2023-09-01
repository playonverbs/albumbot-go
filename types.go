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

type Score uint64

func (s Score) String() string {
	return fmt.Sprintf("%d/30", s)
}

func parseScore(str string) (Score, error) {
	var score string
	_, err := fmt.Sscanf(str, "%d/30", score)
	if err != nil {
		return 0, nil
	}

	val, err := strconv.ParseUint(score, 10, 32)
	if err != nil {
		return 0, err
	}
	return Score(val), nil
}

// TODO: consider just using string for SpotifyURL
type Entry struct {
	ID          int
	Album       string
	DateAdded   time.Time
	SuggestedBy string
	SpotifyURL  url.URL
	Votes       uint
	MeanScore   Score
	Status      ListenStatus
}

func (e Entry) ToValue() []interface{} {
	dateFormat := e.DateAdded.Format(dateLayout)

	return []interface{}{
		e.Album,
		dateFormat,
		e.SuggestedBy,
		e.SpotifyURL.String(),
		e.Votes,
		e.MeanScore.String(),
		e.Status.String(),
	}
}

func NewEntryFromRow(index int, row []interface{}) *Entry {
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

	return &Entry{
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

func CompareDates(d1, d2 time.Time) bool {
	ny, nm, nd := d1.Date()
	ey, em, ed := d2.Date()

	return ny == ey && nm == em && nd == ed
}

// func validateSpotifyLink(string URL)
