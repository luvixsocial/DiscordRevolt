package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

// Event Listener for both Discord and Revolt clients.
func OnEvent(callback func(Event)) {
	if Discord != nil {
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
			eventType := fmt.Sprintf("%T", e)
			callback(Event{
				Name: eventType,
				Type: Message,
				Data: MessageCallback{
					Content: e.Message.Content,
					Author: User{
						ID:       e.Author.ID,
						Username: e.Author.Username,
						Avatar:   e.Author.AvatarURL("128"),
					},
				},
			})
		})

		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.InteractionCreate) {
			eventType := fmt.Sprintf("%T", e)
			callback(Event{
				Name:     eventType,
				Type:     Interaction,
				Platform: "Discord",
				Data: InteractionCallback{
					Name:   e.ApplicationCommandData().Name,
					Fields: convertOptionsToMap(e.ApplicationCommandData().Options),
					Author: User{
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
			eventType := fmt.Sprintf("%T", m)
			authorData, err := Revolt.User(m.Author)

			if err != nil {
				fmt.Println("Failed to get author data:", err)
				return
			}

			callback(Event{
				Name:     eventType,
				Type:     Message,
				Platform: "Revolt",
				Data: MessageCallback{
					Content: m.Content,
					Author: User{
						ID:       authorData.ID,
						Username: authorData.Username,
						Avatar:   authorData.Avatar.URL("128"),
					},
				},
			})
		})
	}
}
