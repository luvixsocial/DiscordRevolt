// Package functions provides core utility helpers for handling
// Discord and Revolt bot operations across platforms, including messaging,
// slash commands, logging, cooldown management, and more.
package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/luvixsocial/whiskercat/types"
	"github.com/sentinelb51/revoltgo"
	"slices"
)

var (
	cooldowns     = make(map[string]time.Time)
	cooldownMutex sync.Mutex
)

func ptr[T any](v T) *T {
	return &v
}

func convertOptionsToMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]string {
	result := make(map[string]string)
	for _, opt := range options {
		if opt != nil && opt.Value != nil {
			result[opt.Name] = fmt.Sprintf("%v", opt.Value)
		}
	}
	return result
}

func convertToDiscordEmbed(embed *types.Embed) *discordgo.MessageEmbed {
	if embed == nil {
		return nil
	}

	em := &discordgo.MessageEmbed{
		Title:       embed.Title,
		Description: embed.Description,
		Color:       embed.Color,
	}

	if embed.URL != nil {
		em.URL = *embed.URL
	}
	if embed.IconURL != nil {
		em.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: *embed.IconURL}
	}
	if embed.PhotoURL != nil {
		em.Image = &discordgo.MessageEmbedImage{URL: *embed.PhotoURL}
	}
	if embed.Footer != nil {
		em.Footer = &discordgo.MessageEmbedFooter{
			Text:    embed.Footer.Text,
			IconURL: embed.Footer.PhotoURL,
		}
	}

	if embed.Fields != nil {
		for _, f := range *embed.Fields {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name:   f.Name,
				Value:  f.Value,
				Inline: f.Inline,
			})
		}
	}
	return em
}

func convertToRevoltEmbed(embed *types.Embed) *revoltgo.MessageEmbed {
	if embed == nil {
		return nil
	}

	description := embed.Description

	// Append fields as Markdown
	if embed.Fields != nil && len(*embed.Fields) > 0 {
		var fieldMarkdown strings.Builder
		fieldMarkdown.WriteString("\n\n")
		for _, f := range *embed.Fields {
			fieldMarkdown.WriteString(fmt.Sprintf("**%s**\n%s\n\n", f.Name, f.Value))
		}
		description += fieldMarkdown.String()
	}

	em := &revoltgo.MessageEmbed{
		Title:       embed.Title,
		Description: description,
		Colour:      fmt.Sprintf("#%06X", embed.Color),
	}

	if embed.URL != nil {
		em.URL = *embed.URL
	}
	if embed.PhotoURL != nil {
		em.Image = &revoltgo.MessageEmbedImage{URL: *embed.PhotoURL}
	}

	return em
}

func Respond(e types.Event, content string, embed *types.Embed, edit *string) (interface{}, error) {
	switch e.Platform {
	case "Discord":
		s := e.Session.(*discordgo.Session)
		switch ctx := e.Context.(type) {
		case *discordgo.MessageCreate:
			if edit != nil {
				em := &discordgo.MessageEdit{ID: *edit, Channel: ctx.ChannelID, Content: &content}
				if embed != nil {
					em.Embeds = &[]*discordgo.MessageEmbed{convertToDiscordEmbed(embed)}
				}
				return s.ChannelMessageEditComplex(em)
			}
			msg := &discordgo.MessageSend{Content: content, Reference: &discordgo.MessageReference{MessageID: ctx.ID, ChannelID: ctx.ChannelID, GuildID: ctx.GuildID}}
			if embed != nil {
				msg.Embeds = []*discordgo.MessageEmbed{convertToDiscordEmbed(embed)}
			}
			return s.ChannelMessageSendComplex(ctx.ChannelID, msg)

		case *discordgo.InteractionCreate:
			if edit != nil {
				ed := &discordgo.WebhookEdit{Content: &content}
				if embed != nil {
					ed.Embeds = &[]*discordgo.MessageEmbed{convertToDiscordEmbed(embed)}
				}
				return s.InteractionResponseEdit(ctx.Interaction, ed)
			}
			r := &discordgo.InteractionResponseData{Content: content}
			if embed != nil {
				r.Embeds = []*discordgo.MessageEmbed{convertToDiscordEmbed(embed)}
			}
			return nil, s.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: r})
		}

	case "Revolt":
		s := e.Session.(*revoltgo.Session)
		var channelID string
		switch ctx := e.Context.(type) {
		case *revoltgo.EventMessage:
			channelID = ctx.Channel
		case *revoltgo.EventMessageUpdate:
			channelID = ctx.Data.Channel
		case *revoltgo.EventMessageReact, *revoltgo.EventMessageRemoveReaction:
			channelID = ctx.(interface{ GetChannelID() string }).GetChannelID()
		default:
			return nil, fmt.Errorf("unsupported Revolt context")
		}
		if edit != nil {
			em := revoltgo.MessageEditData{Content: content}
			if embed != nil {
				em.Embeds = []*revoltgo.MessageEmbed{convertToRevoltEmbed(embed)}
			}
			return s.ChannelMessageEdit(channelID, *edit, em)
		}
		msg := revoltgo.MessageSend{Content: content}
		if embed != nil {
			msg.Embeds = []*revoltgo.MessageEmbed{convertToRevoltEmbed(embed)}
		}
		return s.ChannelMessageSend(channelID, msg)
	}

	return nil, fmt.Errorf("unsupported platform or context")
}

func SendMessage(platform, channelID, content string, embed *types.Embed) (interface{}, error) {
	switch platform {
	case "Discord":
		msg := &discordgo.MessageSend{Content: content}
		if embed != nil {
			msg.Embeds = []*discordgo.MessageEmbed{convertToDiscordEmbed(embed)}
		}
		return Discord.ChannelMessageSendComplex(channelID, msg)

	case "Revolt":
		msg := revoltgo.MessageSend{Content: content}
		if embed != nil {
			msg.Embeds = []*revoltgo.MessageEmbed{convertToRevoltEmbed(embed)}
		}
		return Revolt.ChannelMessageSend(channelID, msg)
	default:
		return nil, fmt.Errorf("unsupported platform")
	}
}

func EditMessage(platform, channelID, messageID, content string, embed *types.Embed) (interface{}, error) {
	switch platform {
	case "Discord":
		em := &discordgo.MessageEdit{ID: messageID, Channel: channelID, Content: &content}
		if embed != nil {
			em.Embeds = &[]*discordgo.MessageEmbed{convertToDiscordEmbed(embed)}
		}
		return Discord.ChannelMessageEditComplex(em)
	case "Revolt":
		em := revoltgo.MessageEditData{Content: content}
		if embed != nil {
			em.Embeds = []*revoltgo.MessageEmbed{convertToRevoltEmbed(embed)}
		}
		return Revolt.ChannelMessageEdit(channelID, messageID, em)
	default:
		return nil, fmt.Errorf("unsupported platform")
	}
}

func DeferInteraction(s *discordgo.Session, interaction *discordgo.Interaction) error {
	return s.InteractionRespond(interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseDeferredChannelMessageWithSource})
}

func GetUserAvatarURL(u types.User) string { return u.Avatar }

func ParseCommand(s string) (string, []string) {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

func LogEvent(e types.Event) { log.Printf("[%s:%s] %+v\n", e.Platform, e.Type, e) }

func IsBotEvent(e types.Event) bool { return e.Bot }

func GetAuthor(e types.Event) types.User {
	switch d := e.Data.(type) {
	case types.MessageCallback:
		return d.Author
	case types.InteractionCallback:
		return d.Author
	}
	return types.User{}
}

func GetChannelID(e types.Event) string {
	switch ctx := e.Context.(type) {
	case *discordgo.MessageCreate:
		return ctx.ChannelID
	case *discordgo.InteractionCreate:
		return ctx.ChannelID
	case *revoltgo.EventMessage:
		return ctx.Channel
	case *revoltgo.EventMessageUpdate:
		return ctx.Data.Channel
	case *revoltgo.EventMessageReact:
		return ctx.ChannelID
	}
	return ""
}

func NewCommand(name, desc string, opts ...*discordgo.ApplicationCommandOption) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: name, Description: desc, Options: opts}
}

func NewOption(name, desc string, typ discordgo.ApplicationCommandOptionType, req bool) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{Name: name, Description: desc, Type: typ, Required: req}
}

func PaginateReply(s *discordgo.Session, cid string, pages []string, d time.Duration) {
	for _, p := range pages {
		s.ChannelMessageSend(cid, p)
		time.Sleep(d)
	}
}

func FindCommandByName(cmds []*discordgo.ApplicationCommand, name string) *discordgo.ApplicationCommand {
	for _, c := range cmds {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func EnsureSlashCommands(s *discordgo.Session, appID, guildID string, cmds []*discordgo.ApplicationCommand) ([]*discordgo.ApplicationCommand, error) {
	existing, err := s.ApplicationCommands(appID, guildID)
	if err != nil {
		return nil, err
	}
	missing := []*discordgo.ApplicationCommand{}
	for _, c := range cmds {
		if FindCommandByName(existing, c.Name) == nil {
			missing = append(missing, c)
		}
	}
	if len(missing) > 0 {
		return s.ApplicationCommandBulkOverwrite(appID, guildID, cmds)
	}
	return existing, nil
}

func GetGuildID(e types.Event) string {
	switch ctx := e.Context.(type) {
	case *discordgo.MessageCreate:
		return ctx.GuildID
	case *discordgo.InteractionCreate:
		return ctx.GuildID
	}
	return ""
}

func GetUsername(e types.Event) string {
	u := GetAuthor(e)
	return fmt.Sprintf("%s (%s)", u.Username, u.ID)
}
func IsAdmin(id string, admins []string) bool { return slices.Contains(admins, id) }

func Cooldown(key string, d time.Duration) bool {
	cooldownMutex.Lock()
	defer cooldownMutex.Unlock()
	if until, exists := cooldowns[key]; exists && time.Now().Before(until) {
		return true
	}
	cooldowns[key] = time.Now().Add(d)
	return false
}

func Retry(fn func() error, attempts int) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return err
}

func Background(fn func()) { go fn() }
