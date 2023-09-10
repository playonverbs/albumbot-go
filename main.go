package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var (
	RemoveCommands *bool
	SpreadsheetID  *string
	GuildID        *string
	DiscordToken   *string
	ReadRange      *string
	Srv            *Service
	ctx            context.Context
)

func init() {
	err := dotenv.Load()
	if err != nil {
		log.Fatalln("Cannot load .env file: ", err)
	}

	RemoveCommands = flag.Bool("rmcmd", true, "remove commands on bot shutdown.")
	SpreadsheetID = flag.String("sheet", "19gDaN_746oECgSAhWlDSRMwhNvuZ70iyjM5l6Q_uL9U", "spreadsheet id to use") // mine
	GuildID = flag.String("guild", os.Getenv("GUILD_ID"), "guild id to use")
	DiscordToken = flag.String("token", os.Getenv("DGU_TOKEN"), "bot access token")
	ReadRange = flag.String("readrange", "Music Club!A2:G", "spreadsheet cell range for IO")

	flag.Parse()
	ctx = context.Background()
}

func init() {
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

	// create google sheets service first to avoid
	Srv, err = NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
}

func main() {
	s, err := discordgo.New("Bot " + *DiscordToken)
	s.SyncEvents = true

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
	log.Println("Opened bot connection...")

	entries, err := Srv.GetSheetEntries(*SpreadsheetID, *ReadRange)
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	fmt.Printf("%#v\n", entries)

	u, err := url.Parse("https://open.spotify.com/album/3HCCUaRSjHSFOe4fqE0BiP?si=roo_D8KGQeiuqqg_2YgCTQ")
	if err != nil {
		log.Fatalf("couldn't parse the URL: %#v", err)
	}
	u.RawQuery = "" // remove share ID

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	log.Println("Press Ctrl-C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down")
}
