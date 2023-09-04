package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func getValues(i *discordgo.Interaction, numFields int) map[string]interface{} {
	var vals = make(map[string]interface{}, numFields)
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
}

var CommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"album-of-the-week": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		values := getValues(i.Interaction, 1)

		fmt.Println(values)

		ents, err := Srv.GetSheetEntries(*SpreadsheetID, *ReadRange)
		if err != nil {
			log.Println("cannot retrieve sheet entries:", err)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					"Album of the Week:\n%s\nVotes: %d\nSubmitted By: %s",
					ents[0].Album, ents[0].Votes, ents[0].SuggestedBy,
				),
			},
		})
	},
	"add-album": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		values := getValues(i.Interaction, 2)
		fmt.Println(values)

		e := NewEntry(values["name"].(string), getMemberName(i.Interaction), values["spotify_link"].(string))
		Srv.AppendSheetEntry(*SpreadsheetID, *ReadRange, e)

		// TODO: add embed/link to show content of what's been added
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Album Added"),
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
}
