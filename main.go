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

	"github.com/FedorLap2006/disgolf"
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

	fmt.Println(*guildID)

	bot, err := disgolf.New(*discordToken)
	if err != nil {
		log.Fatalf("cannot create bot: %#v", err)
	}

	bot.Router.Register(commandAlbumOfTheWeek)
	bot.Router.Register(commandAddAlbum)

	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	bot.AddHandler(bot.Router.HandleInteraction)

	err = bot.Router.Sync(bot.Session, "1146210924395499682", *guildID)
	if err != nil {
		log.Fatal(err)
	}

	err = bot.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	readRange := "Music Club!A2:G"
	entries, err := GetSheetEntries(srv, *SpreadsheetID, readRange)
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	fmt.Printf("%#v", entries)

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
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}
