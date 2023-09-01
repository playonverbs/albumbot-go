package main

import (
	"fmt"
	"time"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
)

func getValues(ctx *disgolf.Ctx) map[string]interface{} {
	var vals map[string]interface{}
	for _, opt := range ctx.Interaction.ApplicationCommandData().Options {
		vals[opt.Name] = opt.Value
	}

	return vals
}

func getMemberName(ctx *disgolf.Ctx) string {
	if ctx.Interaction.Member != nil {
		if nick := ctx.Interaction.Member.Nick; nick != "" {
			return nick
		} else {
			return ctx.Interaction.Member.User.Username
		}
	}

	return ctx.Interaction.User.Username
}

var commandAlbumOfTheWeek = &disgolf.Command{
	Name:        "album-of-the-week",
	Description: "Fetch the current album of the week",
	Handler: disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
		fmt.Println(ctx.Interaction.Member.Nick)

		ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hello World!",
			},
		})
	}),
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "date",
			Description: "date to grab album of the week for, else get the current one",
			Required:    false,
		},
	},
}

var commandAddAlbum = &disgolf.Command{
	Name:        "add-album",
	Description: "Add a new album to the list",
	Handler: disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
		vals := getValues(ctx)
		suggestedBy := getMemberName(ctx)

		fmt.Println(
			vals["name"],
			vals["spotify_link"],
			suggestedBy,
			time.Now().Format(dateLayout),
		)

		ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hello World from AddAlbum",
			},
		})
	}),
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
}

var commandVote = &disgolf.Command{
	Name:        "vote",
	Description: "Vote for an album",
	Handler: disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
		vals := getValues(ctx)
		suggestedBy := getMemberName(ctx)

		fmt.Println(
			vals["name"],
			suggestedBy,
		)
	}),
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "",
			Autocomplete: true,
			Choices:      []*discordgo.ApplicationCommandOptionChoice{},
		},
	},
}
