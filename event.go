package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/luvixsocial/whiskercat/types"
	"github.com/sentinelb51/revoltgo"
)

// OnEvent registers cross-platform event handlers and normalizes them into a common Event format.
func OnEvent(callback func(types.Event)) {
	registerDiscordEvents(callback)
	registerRevoltEvents(callback)
}

func registerDiscordEvents(callback func(types.Event)) {
	if Discord == nil {
		return
	}

	addDiscordHandler := func(eventName string, eventType types.EventType, context interface{}, session *discordgo.Session, bot bool, data interface{}) {
		callback(types.Event{
			Name:     eventName,
			Type:     eventType,
			Platform: "Discord",
			Bot:      bot,
			Context:  context,
			Session:  session,
			Data:     data,
		})
	}

	Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		addDiscordHandler("MessageCreate", types.MessageCreate, e, s, e.Author.Bot, types.MessageCallback{
			Content: e.Content,
			Author:  convertDiscordUser(e.Author),
		})
	})

	Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageUpdate) {
		addDiscordHandler("MessageUpdate", types.MessageUpdate, e, s, e.Author.Bot, types.MessageCallback{
			Content: e.Content,
			Author:  convertDiscordUser(e.Author),
		})
	})

	Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageDelete) {
		addDiscordHandler("MessageDelete", types.MessageDelete, e, s, false, nil)
	})

	Discord.AddHandler(func(s *discordgo.Session, e *discordgo.InteractionCreate) {
		addDiscordHandler("InteractionCreate", types.InteractionCreate, e, s, false, types.InteractionCallback{
			Name:   e.ApplicationCommandData().Name,
			Fields: convertOptionsToMap(e.ApplicationCommandData().Options),
			Data:   e,
			Author: convertDiscordUser(e.Member.User),
		})
	})

	Discord.AddHandler(func(s *discordgo.Session, e *discordgo.TypingStart) {
		addDiscordHandler("TypingStart", types.EventTypingStart, e, s, false, types.User{ID: e.UserID})
	})

	Discord.AddHandler(func(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
		addDiscordHandler("VoiceStateUpdate", types.EventVoiceStateUpdate, e, s, false, nil)
	})

	Discord.AddHandler(func(s *discordgo.Session, e *discordgo.PresenceUpdate) {
		addDiscordHandler("PresenceUpdate", types.EventPresenceUpdate, e, s, false, nil)
	})

	Discord.AddHandler(func(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
		addDiscordHandler("GuildMemberAdd", types.EventGuildMemberAdd, e, s, e.User.Bot, convertDiscordUser(e.User))
	})

	Discord.AddHandler(func(s *discordgo.Session, e *discordgo.GuildMemberRemove) {
		addDiscordHandler("GuildMemberRemove", types.EventGuildMemberRemove, e, s, e.User.Bot, convertDiscordUser(e.User))
	})
}

func registerRevoltEvents(callback func(types.Event)) {
	if Revolt == nil {
		return
	}

	addRevoltHandler := func(eventName string, eventType types.EventType, context interface{}, session *revoltgo.Session, bot bool, data interface{}) {
		callback(types.Event{
			Name:     eventName,
			Type:     eventType,
			Platform: "Revolt",
			Bot:      bot,
			Context:  context,
			Session:  session,
			Data:     data,
		})
	}

	Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventMessage) {
		user, _ := Revolt.User(e.Author)
		addRevoltHandler("MessageCreate", types.MessageCreate, e, s, user.Bot != nil, types.MessageCallback{
			Content: e.Content,
			Author:  convertRevoltUser(user),
		})
	})

	Revolt.AddHandler(func(s *revoltgo.Session, e *revoltgo.EventMessageUpdate) {
		user, _ := Revolt.User(e.Data.Author)
		addRevoltHandler("MessageUpdate", types.MessageUpdate, e, s, user.Bot != nil, types.MessageCallback{
			Content: e.Data.Content,
			Author:  convertRevoltUser(user),
		})
	})

	eventMap := map[string]types.EventType{
		"MessageDelete":  types.MessageDelete,
		"ReactionAdd":    types.ReactionAdd,
		"ReactionRemove": types.ReactionRemove,
		"TypingStart":    types.EventTypingStart,
		"ChannelCreate":  types.EventChannelCreate,
		"ChannelUpdate":  types.EventChannelUpdate,
		"ChannelDelete":  types.EventChannelDelete,
		"UserUpdate":     types.EventUserUpdate,
		"MemberJoin":     types.EventMemberJoin,
		"MemberLeave":    types.EventMemberLeave,
	}

	Revolt.AddHandler(func(s *revoltgo.Session, e interface{}) {
		switch evt := e.(type) {
		case *revoltgo.EventMessageDelete:
			addRevoltHandler("MessageDelete", eventMap["MessageDelete"], evt, s, false, nil)
		case *revoltgo.EventMessageReact:
			addRevoltHandler("ReactionAdd", eventMap["ReactionAdd"], evt, s, false, nil)
		case *revoltgo.EventMessageRemoveReaction:
			addRevoltHandler("ReactionRemove", eventMap["ReactionRemove"], evt, s, false, nil)
		case *revoltgo.EventChannelStartTyping:
			addRevoltHandler("TypingStart", eventMap["TypingStart"], evt, s, false, nil)
		case *revoltgo.EventChannelCreate:
			addRevoltHandler("ChannelCreate", eventMap["ChannelCreate"], evt, s, false, nil)
		case *revoltgo.EventChannelUpdate:
			addRevoltHandler("ChannelUpdate", eventMap["ChannelUpdate"], evt, s, false, nil)
		case *revoltgo.EventChannelDelete:
			addRevoltHandler("ChannelDelete", eventMap["ChannelDelete"], evt, s, false, nil)
		case *revoltgo.EventUserUpdate:
			addRevoltHandler("UserUpdate", eventMap["UserUpdate"], evt, s, false, nil)
		case *revoltgo.EventServerMemberJoin:
			addRevoltHandler("MemberJoin", eventMap["MemberJoin"], evt, s, false, nil)
		case *revoltgo.EventServerMemberLeave:
			addRevoltHandler("MemberLeave", eventMap["MemberLeave"], evt, s, false, nil)
		}
	})
}

func convertDiscordUser(user *discordgo.User) types.User {
	if user == nil {
		return types.User{}
	}
	return types.User{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   user.AvatarURL("128"),
	}
}

func convertRevoltUser(user *revoltgo.User) types.User {
	if user == nil {
		return types.User{}
	}
	return types.User{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   user.Avatar.URL("128"),
	}
}
