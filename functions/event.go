package functions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"

	"whiskercat/types"
)

// OnEvent registers platform-agnostic event listeners and normalizes them.
func OnEvent(callback func(types.Event)) {
	if Discord != nil {
		// Discord Message Events
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
			callback(types.Event{
				Name:     "MessageCreate",
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

		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageUpdate) {
			callback(types.Event{
				Name:     "MessageUpdate",
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

		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageDelete) {
			callback(types.Event{
				Name:     "MessageDelete",
				Type:     types.MessageDelete,
				Platform: "Discord",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		// Discord Interaction Events
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.InteractionCreate) {
			callback(types.Event{
				Name:     "InteractionCreate",
				Type:     types.InteractionCreate,
				Platform: "Discord",
				Bot:      false,
				Context:  e,
				Session:  s,
				Data: types.InteractionCallback{
					Name:   e.ApplicationCommandData().Name,
					Fields: convertOptionsToMap(e.ApplicationCommandData().Options),
					Data:   e,
					Author: types.User{
						ID:       e.Member.User.ID,
						Username: e.Member.User.Username,
						Avatar:   e.Member.User.AvatarURL("128"),
					},
				},
			})
		})

		// Discord Extra Events
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.TypingStart) {
			callback(types.Event{
				Name:     "TypingStart",
				Type:     types.EventTypingStart,
				Platform: "Discord",
				Bot:      false,
				Context:  e,
				Session:  s,
				Data:     types.User{ID: e.UserID},
			})
		})

		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
			callback(types.Event{
				Name:     "VoiceStateUpdate",
				Type:     types.EventVoiceStateUpdate,
				Platform: "Discord",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.PresenceUpdate) {
			callback(types.Event{
				Name:     "PresenceUpdate",
				Type:     types.EventPresenceUpdate,
				Platform: "Discord",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
			callback(types.Event{
				Name:     "GuildMemberAdd",
				Type:     types.EventGuildMemberAdd,
				Platform: "Discord",
				Bot:      e.User.Bot,
				Context:  e,
				Session:  s,
				Data: types.User{
					ID:       e.User.ID,
					Username: e.User.Username,
					Avatar:   e.User.AvatarURL("128"),
				},
			})
		})

		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.GuildMemberRemove) {
			callback(types.Event{
				Name:     "GuildMemberRemove",
				Type:     types.EventGuildMemberRemove,
				Platform: "Discord",
				Bot:      e.User.Bot,
				Context:  e,
				Session:  s,
				Data: types.User{
					ID:       e.User.ID,
					Username: e.User.Username,
					Avatar:   e.User.AvatarURL("128"),
				},
			})
		})
	}

	if Revolt != nil {
		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventMessage) {
			user, _ := Revolt.User(e.Author)
			callback(types.Event{
				Name:     "MessageCreate",
				Type:     types.MessageCreate,
				Platform: "Revolt",
				Bot:      user.Bot != nil,
				Context:  e,
				Session:  s,
				Data: types.MessageCallback{
					Content: e.Content,
					Author: types.User{
						ID:       user.ID,
						Username: user.Username,
						Avatar:   user.Avatar.URL("128"),
					},
				},
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventMessageUpdate) {
			user, _ := Revolt.User(e.Data.Author)
			callback(types.Event{
				Name:     "MessageUpdate",
				Type:     types.MessageUpdate,
				Platform: "Revolt",
				Bot:      user.Bot != nil,
				Context:  e,
				Session:  s,
				Data: types.MessageCallback{
					Content: e.Data.Content,
					Author: types.User{
						ID:       user.ID,
						Username: user.Username,
						Avatar:   user.Avatar.URL("128"),
					},
				},
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventMessageDelete) {
			callback(types.Event{
				Name:     "MessageDelete",
				Type:     types.MessageDelete,
				Platform: "Revolt",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventMessageReact) {
			callback(types.Event{
				Name:     "ReactionAdd",
				Type:     types.ReactionAdd,
				Platform: "Revolt",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventMessageRemoveReaction) {
			callback(types.Event{
				Name:     "ReactionRemove",
				Type:     types.ReactionRemove,
				Platform: "Revolt",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventChannelStartTyping) {
			callback(types.Event{
				Name:     "TypingStart",
				Type:     types.EventTypingStart,
				Platform: "Revolt",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventChannelCreate) {
			callback(types.Event{
				Name:     "ChannelCreate",
				Type:     types.EventChannelCreate,
				Platform: "Revolt",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventChannelUpdate) {
			callback(types.Event{
				Name:     "ChannelUpdate",
				Type:     types.EventChannelUpdate,
				Platform: "Revolt",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventChannelDelete) {
			callback(types.Event{
				Name:     "ChannelDelete",
				Type:     types.EventChannelDelete,
				Platform: "Revolt",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventUserUpdate) {
			callback(types.Event{
				Name:     "UserUpdate",
				Type:     types.EventUserUpdate,
				Platform: "Revolt",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventServerMemberJoin) {
			callback(types.Event{
				Name:     "MemberJoin",
				Type:     types.EventMemberJoin,
				Platform: "Revolt",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})

		Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventServerMemberLeave) {
			callback(types.Event{
				Name:     "MemberLeave",
				Type:     types.EventMemberLeave,
				Platform: "Revolt",
				Bot:      false,
				Context:  e,
				Session:  s,
			})
		})
	}
}
