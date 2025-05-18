package commands

import (
	"fmt"
	"github.com/luvixsocial/whiskercat/types"
)

func TestEmbed(evt types.Event, _ *bool) {
	_, err := Respond(evt, "", &Embed{
		Title:       "Test Embed",
		Description: "This is a test embed.",
		URL:         ptr("https://purrquinox.com/"),
		IconURL:     ptr("https://purrquinox.com/logo.png"),
		Fields: ptr([]EmbedField{
			{
				Name:  "Test Field",
				Value: "This is a test field.",
			},
		}),
		Footer: &EmbedFooter{
			Text:     "This is a test footer.",
			PhotoURL: "https://purrquinox.com/logo.png",
		},
		Color: 0x00FF00,
	}, nil)
	if err != nil {
		fmt.Printf("Error sending embed: %v\n", err)
	}
}
