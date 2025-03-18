package functions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"

	"libdozina/types"
)

// Event Listener for both Discord and Revolt clients.
func OnEvent(callback func(types.Event)) {
	if Discord != nil {
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
			EventType := fmt.Sprintf("%T", e)
			callback(types.Event{
				Name: EventType,
				Type: types.Message,
				Data: types.MessageCallback{
					Content: e.Message.Content,
					Author: types.User{
						ID:       e.Author.ID,
						Username: e.Author.Username,
						Avatar:   e.Author.AvatarURL("128"),
					},
				},
			})
		})

		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.InteractionCreate) {
			EventType := fmt.Sprintf("%T", e)
			callback(types.Event{
				Name:     EventType,
				Type:     types.Interaction,
				Platform: "Discord",
				Data: types.InteractionCallback{
					Name:   e.ApplicationCommandData().Name,
					Fields: convertOptionsToMap(e.ApplicationCommandData().Options),
					Author: types.User{
						ID:       e.Member.User.ID,
						Username: e.Member.User.Username,
						Avatar:   e.Member.User.AvatarURL("128"),
					},
				},
			})
		})
	}

	if Revolt != nil {
		Revolt.AddHandler(func(e *revoltgo.Session, m *revoltgo.EventMessage) {
			EventType := fmt.Sprintf("%T", m)
			authorData, err := Revolt.User(m.Author)

			if err != nil {
				fmt.Println("Failed to get author data:", err)
				return
			}

			callback(types.Event{
				Name:     EventType,
				Type:     types.Message,
				Platform: "Revolt",
				Data: types.MessageCallback{
					Content: m.Content,
					Author: types.User{
						ID:       authorData.ID,
						Username: authorData.Username,
						Avatar:   authorData.Avatar.URL("128"),
					},
				},
			})
		})
	}
}
