package main

import (
	"fmt"
	"time"
	functions "whiskercat/Functions"
	"whiskercat/types"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

// Test
func isolatedtest() {
	functions.Config(&types.Config{
		Discord: types.DiscordConfig{
			ClientID:     "Client ID",
			ClientSecret: "Client Secret",
			Token:        "TOKEN",
		},
		Revolt: types.RevoltConfig{
			Token: "TOKEN",
		},
	})
	functions.Start()
	functions.SetStatus(types.ActivityTypeGame, "Luvix Social", types.Online, nil)

	var stdout bool
	functions.OnEvent(func(evt types.Event) {
		if evt.Bot {
			return
		}

		fmt.Printf("Received event: %v", evt)

		if evt.Type == "Message" {
			if data, ok := evt.Data.(types.MessageCallback); ok && data.Content == "dev:start" {
				stdout = true

				_, err := functions.Respond(evt, "Enabled developer mode.", nil, nil)
				if err != nil {
					fmt.Printf("Error sending message: %v\n", err.Error())
				}
				return
			} else if data, ok := evt.Data.(types.MessageCallback); ok && data.Content == "dev:stop" {
				stdout = false

				_, err := functions.Respond(evt, "Disabled developer mode.", nil, nil)
				if err != nil {
					fmt.Printf("Error sending message: %v\n", err.Error())
				}
				return
			}

			if stdout {
				_, err := functions.Respond(evt, "", &types.Embed{
					Title:       "Event Recieved.",
					Description: fmt.Sprintf("%+v", evt),
					Color:       0x00FF00,
				}, nil)
				if err != nil {
					fmt.Printf("Error sending message: %v\n", err.Error())
				}
			}

			// Commands
			if a, ok := evt.Data.(types.MessageCallback); ok {
				if a.Content == "ping" {
					start := time.Now()

					// Send initial message
					msgRaw, err := functions.Respond(evt, "Pinging...", nil, nil)
					if err != nil {
						fmt.Printf("Error sending message: %v\n", err.Error())
						return
					}

					latency := time.Since(start).Milliseconds()
					pingContent := fmt.Sprintf("🏓 Pong! %dms", latency)

					switch msg := msgRaw.(type) {
					case *discordgo.Message:
						_, err = functions.Respond(evt, pingContent, nil, &msg.ID)
					case *revoltgo.Message:
						_, err = functions.Respond(evt, pingContent, nil, &msg.ID)
					}

					if err != nil {
						fmt.Printf("Error editing message: %v\n", err.Error())
					}
				} else if a.Content == "test" {
					_, err := functions.Respond(evt, "Received test event!", nil, nil)
					if err != nil {
						fmt.Printf("Error sending message: %v\n", err.Error())
					}
				}
			}
		} else if evt.Type == "MessageUpdate" {
			if stdout {
				var channelID string

				switch evt.Platform {
				case "Discord":
					if ctx, ok := evt.Context.(*discordgo.MessageUpdate); ok {
						channelID = ctx.ChannelID
					}
				case "Revolt":
					if ctx, ok := evt.Context.(*revoltgo.EventMessageUpdate); ok {
						channelID = ctx.Channel
					}
				}

				if channelID != "" {
					_, err := functions.SendMessage(evt.Platform, channelID, "", &types.Embed{
						Title:       "Event Received.",
						Description: fmt.Sprintf("%+v", evt),
						Color:       0x00FF00,
					})
					if err != nil {
						fmt.Printf("Error sending message: %v\n", err.Error())
					}
				}
			}
		} else if evt.Type == "ReactionAdd" {
			if stdout {
				var channelID string

				switch evt.Platform {
				case "Discord":
					if ctx, ok := evt.Context.(*discordgo.MessageReactionAdd); ok {
						channelID = ctx.ChannelID
					}
				case "Revolt":
					if ctx, ok := evt.Context.(*revoltgo.EventMessageReact); ok {
						channelID = ctx.ChannelID
					}
				}

				if channelID != "" {
					_, err := functions.SendMessage(evt.Platform, channelID, "", &types.Embed{
						Title:       "Event Received.",
						Description: fmt.Sprintf("%+v", evt),
						Color:       0x00FF00,
					})
					if err != nil {
						fmt.Printf("Error sending message: %v\n", err.Error())
					}
				}
			}
		} else if evt.Type == "ReactionRemove" {
			if stdout {
				var channelID string

				switch evt.Platform {
				case "Discord":
					if ctx, ok := evt.Context.(*discordgo.MessageReactionRemove); ok {
						channelID = ctx.ChannelID
					}
				case "Revolt":
					if ctx, ok := evt.Context.(*revoltgo.EventMessageRemoveReaction); ok {
						channelID = ctx.ChannelID
					}
				}

				if channelID != "" {
					_, err := functions.SendMessage(evt.Platform, channelID, "", &types.Embed{
						Title:       "Event Received.",
						Description: fmt.Sprintf("%+v", evt),
						Color:       0x00FF00,
					})
					if err != nil {
						fmt.Printf("Error sending message: %v\n", err.Error())
					}
				}
			}
		}
	})

	select {}
}
