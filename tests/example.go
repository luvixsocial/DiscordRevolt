package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

func example() {
	Config(&AuthConfig{
		Discord: DiscordConfig{
			ClientID:     "YOUR_CLIENT_ID",
			ClientSecret: "YOUR_CLIENT_SECRET",
			Token:        "YOUR_DISCORD_TOKEN",
		},
		Revolt: RevoltConfig{
			Token: "YOUR_REVOLT_TOKEN",
		},
	})

	dSession, _ := Start()
	SetStatus(ActivityTypeGame, "Luvix Social", Online, nil)

	var stdout bool

	OnEvent(func(evt Event) {
		if evt.Bot {
			return
		}

		fmt.Printf("Received event: %v\n", evt)

		switch evt.Type {
		case MessageCreate:
			handleCommand(evt, &stdout, getMessageContent(evt))
		case InteractionCreate:
			if evt.Platform == "Discord" {
				interaction := evt.Context.(*discordgo.InteractionCreate)
				if interaction.Type == discordgo.InteractionApplicationCommand {
					handleCommand(evt, &stdout, interaction.ApplicationCommandData().Name)
				}
			}
		case MessageUpdate, MessageDelete,
			ReactionAdd, ReactionRemove,
			EventTypingStart, EventVoiceStateUpdate,
			EventPresenceUpdate, EventGuildMemberAdd,
			EventGuildMemberRemove, EventChannelCreate,
			EventChannelUpdate, EventChannelDelete,
			EventUserUpdate, EventMemberJoin,
			EventMemberLeave:
			if stdout {
				emitEventEmbed(evt, getChannelID(evt))
			}
		}
	})

	commands := []*discordgo.ApplicationCommand{
		NewCommand("ping", "Check latency"),
		NewCommand("enable_dev", "Enable developer mode"),
		NewCommand("disable_dev", "Disable developer mode"),
		NewCommand("test", "Test event reply"),
		NewCommand("test_embed", "Send a test embed"),
	}
	EnsureSlashCommands(dSession, dSession.State.Application.ID, "", commands)

	select {}
}

func handleCommand(evt Event, stdout *bool, command string) {
	switch command {
	case "enable_dev":
		*stdout = true
		Respond(evt, "Enabled developer mode.", nil, nil)
	case "disable_dev":
		*stdout = false
		Respond(evt, "Disabled developer mode.", nil, nil)
	case "ping":
		start := time.Now()
		msgRaw, err := Respond(evt, "Pinging...", nil, nil)
		if err != nil {
			fmt.Printf("Error sending ping: %v\n", err)
			return
		}
		latency := time.Since(start).Milliseconds()
		pong := fmt.Sprintf("Pong! %dms", latency)
		switch msg := msgRaw.(type) {
		case *discordgo.Message:
			Respond(evt, pong, nil, &msg.ID)
		case *revoltgo.Message:
			Respond(evt, pong, nil, &msg.ID)
		}
	case "test":
		Respond(evt, "Received test event!", nil, nil)
	case "test_embed":
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
	default:
		if *stdout {
			Respond(evt, "", &Embed{
				Title:       "Event Received",
				Description: fmt.Sprintf("%+v", evt),
				Color:       0x00FF00,
			}, nil)
		}
	}
}

func getMessageContent(evt Event) string {
	if data, ok := evt.Data.(MessageCallback); ok {
		return data.Content
	}
	return ""
}

func emitEventEmbed(evt Event, channelID string) {
	if channelID == "" {
		return
	}
	_, err := SendMessage(evt.Platform, channelID, "", &Embed{
		Title:       "Event Received",
		Description: fmt.Sprintf("%+v", evt),
		Color:       0x00FF00,
	})
	if err != nil {
		fmt.Printf("Error sending embed: %v\n", err)
	}
}

func getChannelID(evt Event) string {
	switch evt.Platform {
	case "Discord":
		switch ctx := evt.Context.(type) {
		case *discordgo.MessageUpdate:
			return ctx.ChannelID
		case *discordgo.MessageDelete:
			return ctx.ChannelID
		case *discordgo.MessageReactionAdd:
			return ctx.ChannelID
		case *discordgo.MessageReactionRemove:
			return ctx.ChannelID
		case *discordgo.TypingStart:
			return ctx.ChannelID
		case *discordgo.VoiceStateUpdate:
			return ctx.ChannelID
		case *discordgo.PresenceUpdate:
			return ctx.User.ID
		case *discordgo.GuildMemberAdd:
			return ctx.GuildID
		case *discordgo.GuildMemberRemove:
			return ctx.GuildID
		}
	case "Revolt":
		switch ctx := evt.Context.(type) {
		case *revoltgo.EventMessageUpdate:
			return ctx.Channel
		case *revoltgo.EventMessageDelete:
			return ctx.Channel
		case *revoltgo.EventMessageReact:
			return ctx.ChannelID
		case *revoltgo.EventMessageRemoveReaction:
			return ctx.ChannelID
		case *revoltgo.EventChannelStartTyping:
			return ctx.ID
		case *revoltgo.EventChannelCreate:
			return ctx.Channel.ID
		case *revoltgo.EventChannelUpdate:
			return ctx.ID
		case *revoltgo.EventChannelDelete:
			return ctx.ID
		case *revoltgo.EventServerMemberJoin:
			return ctx.ID
		case *revoltgo.EventServerMemberLeave:
			return ctx.ID
		case *revoltgo.EventUserUpdate:
			return ctx.ID
		}
	}
	return ""
}
