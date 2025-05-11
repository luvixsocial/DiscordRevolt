package functions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"

	"whiskercat/types"
)

// Event Listener for both Discord and Revolt clients.
func OnEvent(callback func(types.Event)) {
	if Discord != nil {
		// MessageCreate
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
			EventType := fmt.Sprintf("%T", e)
			callback(types.Event{
				Name:     EventType,
				Type:     types.Message,
				Platform: "Discord",
				Bot:      e.Author.Bot,
				Context:  e,
				Session:  s,
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

		// MessageUpdate
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageUpdate) {
			EventType := fmt.Sprintf("%T", e)
			callback(types.Event{
				Name:     EventType,
				Type:     types.MessageUpdate,
				Platform: "Discord",
				Bot:      e.Author.Bot,
				Context:  e,
				Session:  s,
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

		// MessageReactionAdd
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
			EventType := fmt.Sprintf("%T", e)
			callback(types.Event{
				Name:     EventType,
				Type:     types.ReactionAdd,
				Platform: "Discord",
				Bot:      e.Member.User.Bot,
				Context:  e,
				Session:  s,
				Data:     nil,
			})
		})

		// MessageReactionRemove
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
			EventType := fmt.Sprintf("%T", e)
			callback(types.Event{
				Name:     EventType,
				Type:     types.ReactionRemove,
				Platform: "Discord",
				Bot:      false,
				Context:  e,
				Session:  s,
				Data:     nil,
			})
		})

		// InteractionCreate
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.InteractionCreate) {
			EventType := fmt.Sprintf("%T", e)
			callback(types.Event{
				Name:     EventType,
				Type:     types.Interaction,
				Platform: "Discord",
				Bot:      false,
				Context:  e,
				Session:  s,
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
		// MessageCreate
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
				Bot:      authorData.Bot != nil,
				Context:  m,
				Session:  e,
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

		// MessageUpdate
		Revolt.AddHandler(func(e *revoltgo.Session, m *revoltgo.EventMessageUpdate) {
			EventType := fmt.Sprintf("%T", m)
			authorData, err := Revolt.User(m.Data.Author)

			if err != nil {
				fmt.Println("Failed to get author data:", err)
				return
			}

			callback(types.Event{
				Name:     EventType,
				Type:     types.MessageUpdate,
				Platform: "Revolt",
				Bot:      authorData.Bot != nil,
				Context:  m,
				Session:  e,
				Data: types.MessageCallback{
					Content: m.Data.Content,
					Author: types.User{
						ID:       authorData.ID,
						Username: authorData.Username,
						Avatar:   authorData.Avatar.URL("128"),
					},
				},
			})
		})

		// MessageReactionAdd
		Revolt.AddHandler(func(e *revoltgo.Session, m *revoltgo.EventMessageReact) {
			EventType := fmt.Sprintf("%T", m)

			callback(types.Event{
				Name:     EventType,
				Type:     types.ReactionAdd,
				Platform: "Revolt",
				Bot:      false,
				Context:  m,
				Session:  e,
				Data:     nil,
			})
		})

		// MessageReactionRemove
		Revolt.AddHandler(func(e *revoltgo.Session, m *revoltgo.EventMessageRemoveReaction) {
			EventType := fmt.Sprintf("%T", m)

			callback(types.Event{
				Name:     EventType,
				Type:     types.ReactionRemove,
				Platform: "Revolt",
				Bot:      false,
				Context:  m,
				Session:  e,
				Data:     nil,
			})
		})
	}
}
