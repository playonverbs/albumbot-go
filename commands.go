package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func getValues(i *discordgo.Interaction) map[string]interface{} {
	var vals map[string]interface{}
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
		},
	},
}

var CommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"album-of-the-week": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		values := getValues(i.Interaction)

		fmt.Println(values)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there from album-of-the-week",
			},
		})
	},
	"add-album": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		values := getValues(i.Interaction)
		fmt.Println(values)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Hey there from add-album"),
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
