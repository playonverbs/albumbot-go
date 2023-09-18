package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
)

func getValues(i *discordgo.Interaction) map[string]interface{} {
	var vals = make(map[string]interface{})
	for _, opt := range i.ApplicationCommandData().Options {
		vals[opt.Name] = opt.Value
	}

	return vals
}

func getMemberName(i *discordgo.Interaction) string {
	if i.Member != nil {
		if nick := i.Member.Nick; nick != "" {
			return nick
		} else {
			return i.Member.User.Username
		}
	}

	return i.User.Username
}

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "hi",
		Description: "be polite",
	},
	{
		Name:        "album-of-the-week",
		Description: "Fetch the current album of the week",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "date to grab album of the week for, else get the current one",
				Required:    false,
			},
		},
	},
	{
		Name:        "add-album",
		Description: "Add a new album to the list",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "name of the album to add",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "spotify_link",
				Description: "spotify link to the album",
				Required:    true,
			},
		},
	},
	{
		Name:        "vote",
		Description: "Vote for an album",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "album-name",
				Description:  "Name of the album to vote for",
				Autocomplete: true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "The-Kickback",
						Value: "The Kickback",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "score",
				Description: "Score to give the album",
				MaxValue:    30.0,
			},
		},
	},
	{
		Name:        "roll",
		Description: "Randomly pick from either albums not-listened to (default), previously listened to, or all albums.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "filter",
				Description:  "optional condition to filter roll by",
				Autocomplete: true,
				Required:     false,
			},
		},
	},
}

var CommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"hi": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "hi!"},
		})
	},
	"album-of-the-week": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		values := getValues(i.Interaction)

		fmt.Println(values)

		ents, err := Srv.GetSheetEntries(*SpreadsheetID, *ReadRange)
		if err != nil {
			log.Println("cannot retrieve sheet entries:", err)
		}

		weeksEnts := ents.AlbumsInWeek(time.Now())
		sort.Sort(sort.Reverse(ByVotes{weeksEnts}))

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					"Album of the Week: %s\nVotes: %d\nSubmitted By: %s\nLink: %s",
					weeksEnts[0].Album,
					weeksEnts[0].Votes,
					weeksEnts[0].SuggestedBy,
					weeksEnts[0].SpotifyURL.String(),
				),
			},
		})
	},
	"add-album": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		values := getValues(i.Interaction)
		fmt.Println(values)

		e := NewAlbum(values["name"].(string), getMemberName(i.Interaction), values["spotify_link"].(string))
		Srv.AppendSheetEntry(*SpreadsheetID, *ReadRange, e)

		// TODO: add embed/link to show content of what's been added
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Album Added: %s", e.Album),
			},
		})
	},
	"vote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there from vote",
			},
		})
	},
	"roll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// values := getValues(i.Interaction)

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			ents, err := Srv.GetSheetEntries(*SpreadsheetID, *ReadRange)
			if err != nil {
				log.Println(err)
			}

			r := ents.Rand()

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						"Picked: %s: %s",
						r.Album,
						r.SpotifyURL.String(),
					),
				},
			})
		case discordgo.InteractionApplicationCommandAutocomplete:
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: commandChoices["roll"],
				},
			})
		}
	},
}

// autocomplete choices for each command
var commandChoices = map[string][]*discordgo.ApplicationCommandOptionChoice{
	"roll": {
		{
			Name:  "Not Listened",
			Value: NotListened.String(),
		},
		{
			Name:  "Listened",
			Value: Listened.String(),
		},
		{
			Name:  "All",
			Value: "all",
		},
	},
}
