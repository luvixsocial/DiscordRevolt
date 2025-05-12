package tests

import (
	"fmt"
	"time"

	"whiskercat/functions"
	"whiskercat/types"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

func isolated_test() {
	// Initialize bot configuration for both platforms
	functions.Config(&types.Config{
		Discord: types.DiscordConfig{
			ClientID:     "YOUR_DISCORD_CLIENT_ID",
			ClientSecret: "YOUR_DISCORD_CLIENT_SECRET",
			Token:        "YOUR_DISCORD_TOKEN",
		},
		Revolt: types.RevoltConfig{
			Token: "YOUR_REVOLT_TOKEN",
		},
	})

	// Start the clients
	functions.Start()

	// Set initial bot status
	functions.SetStatus(types.ActivityTypeGame, "Luvix Social", types.Online, nil)

	// Enable live debug output toggle
	var stdout bool

	// Event listener
	functions.OnEvent(func(evt types.Event) {
		if evt.Bot {
			return // Ignore bot messages
		}

		fmt.Printf("Received event: %v\n", evt)

		switch evt.Type {
		case types.MessageCreate:
			handleMessageCreate(evt, &stdout)
		case types.MessageUpdate:
			if stdout {
				emitEventEmbed(evt, getChannelID(evt))
			}
		case types.MessageDelete:
			if stdout {
				emitEventEmbed(evt, getChannelID(evt))
			}
		case types.ReactionAdd, types.ReactionRemove:
			if stdout {
				emitEventEmbed(evt, getChannelID(evt))
			}
		}
	})

	select {} // Block forever
}

// handleMessageCreate processes MessageCreate events, including toggle and command logic.
func handleMessageCreate(evt types.Event, stdout *bool) {
	data, ok := evt.Data.(types.MessageCallback)
	if !ok {
		return
	}

	switch data.Content {
	case "dev:start":
		*stdout = true
		functions.Respond(evt, "Enabled developer mode.", nil, nil)

	case "dev:stop":
		*stdout = false
		functions.Respond(evt, "Disabled developer mode.", nil, nil)

	case "ping":
		start := time.Now()
		msgRaw, err := functions.Respond(evt, "Pinging...", nil, nil)
		if err != nil {
			fmt.Printf("Error sending ping: %v\n", err)
			return
		}

		latency := time.Since(start).Milliseconds()
		pong := fmt.Sprintf("🏓 Pong! %dms", latency)

		switch msg := msgRaw.(type) {
		case *discordgo.Message:
			functions.Respond(evt, pong, nil, &msg.ID)
		case *revoltgo.Message:
			functions.Respond(evt, pong, nil, &msg.ID)
		}

	case "test":
		functions.Respond(evt, "Received test event!", nil, nil)

	case "test:embed":
		_, err := functions.Respond(evt, "", &types.Embed{
			Title:       "Test Embed",
			Description: "This is a test embed.",
			Color:       0x00FF00,
		}, nil)
		if err != nil {
			fmt.Printf("Error sending embed: %v\n", err)
		}

	default:
		if *stdout {
			functions.Respond(evt, "", &types.Embed{
				Title:       "Event Received",
				Description: fmt.Sprintf("%+v", evt),
				Color:       0x00FF00,
			}, nil)
		}
	}
}

// emitEventEmbed sends a green embed with event data to the provided channel.
func emitEventEmbed(evt types.Event, channelID string) {
	if channelID == "" {
		return
	}
	_, err := functions.SendMessage(evt.Platform, channelID, "", &types.Embed{
		Title:       "Event Received",
		Description: fmt.Sprintf("%+v", evt),
		Color:       0x00FF00,
	})
	if err != nil {
		fmt.Printf("Error sending embed: %v\n", err)
	}
}

// getChannelID extracts the channel ID from the event's context.
func getChannelID(evt types.Event) string {
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
		}
	}
	return ""
}
