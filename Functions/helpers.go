package functions

import (
	"fmt"
	"whiskercat/types"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

// Helper function to convert []*discordgo.ApplicationCommandInteractionDataOption to map[string]string
func convertOptionsToMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]string {
	fields := make(map[string]string)
	for _, option := range options {
		if option != nil && option.Value != nil {
			fields[option.Name] = fmt.Sprintf("%v", option.Value)
		}
	}
	return fields
}

// Helper function to respond to any event.
func Respond(e types.Event, content string, embed *types.Embed, edit *string) (interface{}, error) {
	switch e.Platform {
	case "Discord":
		session := e.Session.(*discordgo.Session)

		switch ctx := e.Context.(type) {
		case *discordgo.MessageCreate:
			if edit != nil {
				// Edit message
				editMsg := &discordgo.MessageEdit{
					ID:      *edit,
					Channel: ctx.ChannelID,
					Content: &content,
				}
				if embed != nil {
					embeds := []*discordgo.MessageEmbed{{
						Title:       embed.Title,
						Description: embed.Description,
						Color:       embed.Color,
					}}
					editMsg.Embeds = &embeds
				}
				return session.ChannelMessageEditComplex(editMsg)
			}

			// Send message
			msg := &discordgo.MessageSend{
				Content: content,
				Reference: &discordgo.MessageReference{
					MessageID: ctx.ID,
					ChannelID: ctx.ChannelID,
					GuildID:   ctx.GuildID,
				},
			}
			if embed != nil {
				msg.Embeds = []*discordgo.MessageEmbed{{
					Title:       embed.Title,
					Description: embed.Description,
					Color:       embed.Color,
				}}
			}
			return session.ChannelMessageSendComplex(ctx.ChannelID, msg)

		case *discordgo.InteractionCreate:
			if edit != nil {
				// Edit interaction response
				editData := &discordgo.WebhookEdit{
					Content: &content,
				}
				if embed != nil {
					embeds := []*discordgo.MessageEmbed{{
						Title:       embed.Title,
						Description: embed.Description,
						Color:       embed.Color,
					}}
					editData.Embeds = &embeds
				}
				resp, err := session.InteractionResponseEdit(ctx.Interaction, editData)
				return resp, err
			}

			// Respond to interaction
			data := &discordgo.InteractionResponseData{
				Content: content,
			}
			if embed != nil {
				data.Embeds = []*discordgo.MessageEmbed{{
					Title:       embed.Title,
					Description: embed.Description,
					Color:       embed.Color,
				}}
			}
			err := session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: data,
			})
			return nil, err
		}

	case "Revolt":
		session := e.Session.(*revoltgo.Session)
		ctx, ok := e.Context.(*revoltgo.EventMessage)
		if !ok {
			break
		}

		if edit != nil {
			editMsg := revoltgo.MessageEditData{
				Content: content,
			}
			if embed != nil {
				editMsg.Embeds = []*revoltgo.MessageEmbed{{
					Title:       embed.Title,
					Description: embed.Description,
					Colour:      fmt.Sprintf("#%06X", embed.Color),
				}}
			}
			return session.ChannelMessageEdit(ctx.Channel, *edit, editMsg)
		}

		sendMsg := revoltgo.MessageSend{
			Content: content,
		}
		if embed != nil {
			sendMsg.Embeds = []*revoltgo.MessageEmbed{{
				Title:       embed.Title,
				Description: embed.Description,
				Colour:      fmt.Sprintf("#%06X", embed.Color),
			}}
		}
		return session.ChannelMessageSend(ctx.Channel, sendMsg)
	}

	return nil, fmt.Errorf("unsupported platform or context")
}

// Helper function to send a message to any channel.
func SendMessage(platform string, channelID string, content string, embed *types.Embed) (interface{}, error) {
	switch platform {
	case "Discord":
		msg := &discordgo.MessageSend{
			Content: content,
		}
		if embed != nil {
			msg.Embeds = []*discordgo.MessageEmbed{{
				Title:       embed.Title,
				Description: embed.Description,
				Color:       embed.Color,
			}}
		}
		return Discord.ChannelMessageSendComplex(channelID, msg)

	case "Revolt":
		sendMsg := revoltgo.MessageSend{
			Content: content,
		}
		if embed != nil {
			sendMsg.Embeds = []*revoltgo.MessageEmbed{{
				Title:       embed.Title,
				Description: embed.Description,
				Colour:      fmt.Sprintf("#%06X", embed.Color),
			}}
		}
		return Revolt.ChannelMessageSend(channelID, sendMsg)

	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}
