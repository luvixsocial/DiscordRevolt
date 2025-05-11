package functions

import (
	"fmt"
	"whiskercat/types"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

// convertOptionsToMap transforms Discord command options into a simple map[string]string.
//
// Used primarily for parsing slash command fields.
func convertOptionsToMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]string {
	fields := make(map[string]string)
	for _, option := range options {
		if option != nil && option.Value != nil {
			fields[option.Name] = fmt.Sprintf("%v", option.Value)
		}
	}
	return fields
}

// Respond sends or edits a message or interaction response based on the event platform and context.
//
// - `e`: Event to respond to
// - `content`: Message content to send
// - `embed`: Optional embed to include (can be nil)
// - `edit`: Optional message ID to edit (nil = create new)
//
// Returns the raw response (message or interaction) and an error, if any.
func Respond(e types.Event, content string, embed *types.Embed, edit *string) (interface{}, error) {
	switch e.Platform {
	case "Discord":
		session := e.Session.(*discordgo.Session)

		switch ctx := e.Context.(type) {

		// Respond to a message (new or edited)
		case *discordgo.MessageCreate:
			if edit != nil {
				// Edit existing message
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

			// Send new message
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

		// Respond to interaction (slash command, etc.)
		case *discordgo.InteractionCreate:
			if edit != nil {
				// Edit previous interaction response
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

			// Send initial interaction response
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
			// Edit existing message
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

		// Send new message
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

// SendMessage sends a message directly to a given channel on a supported platform.
//
// - `platform`: "Discord" or "Revolt"
// - `channelID`: Channel ID where the message will be sent
// - `content`: Plain text message
// - `embed`: Optional embed to include
//
// Returns the sent message and error, if any.
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
