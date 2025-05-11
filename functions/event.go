package functions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"

	"whiskercat/types"
)

// OnEvent registers a unified event listener for both Discord and Revolt platforms.
// It accepts a callback function that will be invoked with a normalized `types.Event` object
// regardless of which platform the original event originated from.
//
// This helps abstract away platform-specific differences and allows uniform event handling.
func OnEvent(callback func(types.Event)) {
	// Discord event listeners
	if Discord != nil {
		// Handle Discord message creation
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
			callback(types.Event{
				Name:     fmt.Sprintf("%T", e),
				Type:     types.MessageCreate,
				Platform: "Discord",
				Bot:      e.Author.Bot,
				Context:  e,
				Session:  s,
				Data: types.MessageCallback{
					Content: e.Content,
					Author: types.User{
						ID:       e.Author.ID,
						Username: e.Author.Username,
						Avatar:   e.Author.AvatarURL("128"),
					},
				},
			})
		})

		// Handle Discord message updates
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageUpdate) {
			callback(types.Event{
				Name:     fmt.Sprintf("%T", e),
				Type:     types.MessageUpdate,
				Platform: "Discord",
				Bot:      e.Author.Bot,
				Context:  e,
				Session:  s,
				Data: types.MessageCallback{
					Content: e.Content,
					Author: types.User{
						ID:       e.Author.ID,
						Username: e.Author.Username,
						Avatar:   e.Author.AvatarURL("128"),
					},
				},
			})
		})

		// Handle Discord reaction additions
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
			callback(types.Event{
				Name:     fmt.Sprintf("%T", e),
				Type:     types.ReactionAdd,
				Platform: "Discord",
				Bot:      e.Member.User.Bot,
				Context:  e,
				Session:  s,
				Data:     nil,
			})
		})

		// Handle Discord reaction removals
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
			callback(types.Event{
				Name:     fmt.Sprintf("%T", e),
				Type:     types.ReactionRemove,
				Platform: "Discord",
				Bot:      false,
				Context:  e,
				Session:  s,
				Data:     nil,
			})
		})

		// Handle Discord interactions (e.g., slash commands)
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.InteractionCreate) {
			callback(types.Event{
				Name:     fmt.Sprintf("%T", e),
				Type:     types.InteractionCreate,
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

	// Revolt event listeners
	if Revolt != nil {
		// Handle Revolt message creation
		Revolt.AddHandler(func(s *revoltgo.Session, m *revoltgo.EventMessage) {
			authorData, err := Revolt.User(m.Author)
			if err != nil {
				fmt.Println("Failed to fetch Revolt user data:", err)
				return
			}

			callback(types.Event{
				Name:     fmt.Sprintf("%T", m),
				Type:     types.MessageCreate,
				Platform: "Revolt",
				Bot:      authorData.Bot != nil,
				Context:  m,
				Session:  s,
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

		// Handle Revolt message updates
		Revolt.AddHandler(func(s *revoltgo.Session, m *revoltgo.EventMessageUpdate) {
			authorData, err := Revolt.User(m.Data.Author)
			if err != nil {
				fmt.Println("Failed to fetch Revolt user data:", err)
				return
			}

			callback(types.Event{
				Name:     fmt.Sprintf("%T", m),
				Type:     types.MessageUpdate,
				Platform: "Revolt",
				Bot:      authorData.Bot != nil,
				Context:  m,
				Session:  s,
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

		// Handle Revolt reaction additions
		Revolt.AddHandler(func(s *revoltgo.Session, m *revoltgo.EventMessageReact) {
			callback(types.Event{
				Name:     fmt.Sprintf("%T", m),
				Type:     types.ReactionAdd,
				Platform: "Revolt",
				Bot:      false,
				Context:  m,
				Session:  s,
				Data:     nil,
			})
		})

		// Handle Revolt reaction removals
		Revolt.AddHandler(func(s *revoltgo.Session, m *revoltgo.EventMessageRemoveReaction) {
			callback(types.Event{
				Name:     fmt.Sprintf("%T", m),
				Type:     types.ReactionRemove,
				Platform: "Revolt",
				Bot:      false,
				Context:  m,
				Session:  s,
				Data:     nil,
			})
		})
	}
}
