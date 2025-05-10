package functions

import (
	"fmt"
	"libdozina/types"

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
func Respond(e types.Event, content string) error {
	switch e.Platform {
	case "Discord":
		switch ctx := e.Context.(type) {
		case *discordgo.MessageCreate:
			session := e.Session.(*discordgo.Session)
			_, err := session.ChannelMessageSendReply(ctx.ChannelID, content, &discordgo.MessageReference{
				MessageID: ctx.ID,
				ChannelID: ctx.ChannelID,
				GuildID:   ctx.GuildID,
			})
			return err

		case *discordgo.InteractionCreate:
			session := e.Session.(*discordgo.Session)
			err := session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			return err
		}

	case "Revolt":
		if ctx, ok := e.Context.(*revoltgo.EventMessage); ok {
			session := e.Session.(*revoltgo.Session)
			_, err := session.ChannelMessageSend(ctx.Channel, revoltgo.MessageSend{
				Content: content,
			})
			return err
		}
	}

	return fmt.Errorf("unsupported platform or context")
}
