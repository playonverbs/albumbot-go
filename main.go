package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	ctx := context.Background()
	err := dotenv.Load()

	discordToken := flag.String("token", os.Getenv("DGU_TOKEN"), "bot access token")
	guildID := flag.String("guild", os.Getenv("GUILD_ID"), "guild id to use")
	// SpreadsheetID := flag.String("sheet", "1QYP9sN_VeDDzbiNI03IRPzAR1KafFv5YjUTVWq6fh_k", "spreadsheet id to use") // nicks
	SpreadsheetID := flag.String("sheet", "19gDaN_746oECgSAhWlDSRMwhNvuZ70iyjM5l6Q_uL9U", "spreadsheet id to use") // mine

	flag.Parse()

	s, err := discordgo.New("Bot " + *discordToken)

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandler[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	defer s.Close()

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := GetClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	readRange := "Music Club!A2:G"
	entries, err := GetSheetEntries(srv, *SpreadsheetID, readRange)
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	fmt.Printf("%#v\n", entries)

	u, err := url.Parse("https://open.spotify.com/album/3HCCUaRSjHSFOe4fqE0BiP?si=roo_D8KGQeiuqqg_2YgCTQ")
	if err != nil {
		log.Fatalf("couldn't parse the URL: %#v", err)
	}
	u.RawQuery = "" // remove share ID

	e := &Entry{
		Album:       "The Kickback",
		DateAdded:   time.Now(),
		SuggestedBy: "Niam",
		SpotifyURL:  *u,
		Votes:       3,
		MeanScore:   30,
		Status:      NotListened,
	}

	err = AppendSheetEntry(srv, *SpreadsheetID, readRange, e)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	log.Println("Press Ctrl-C to exit")
	<-stop

	log.Println("Gracefully shutting down")
}
