// Package functions provides core utility helpers for handling
// Discord and Revolt bot operations across platforms, including messaging,
// slash commands, logging, cooldown management, and more.
package functions

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"whiskercat/types"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
	"slices"
)

// Global cooldown map to track command usage per key
var cooldowns = make(map[string]time.Time)
var cooldownLock sync.Mutex

// convertOptionsToMap transforms Discord command options into a map[string]string
// to simplify the parsing and use of command arguments.
func convertOptionsToMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]string {
	fields := make(map[string]string)
	for _, option := range options {
		if option != nil && option.Value != nil {
			fields[option.Name] = fmt.Sprintf("%v", option.Value)
		}
	}
	return fields
}

// Respond handles unified response logic for Discord and Revolt events.
// It detects whether to send a new message or edit an existing one.
//
// Parameters:
// - e: types.Event struct containing platform, context, and session.
// - content: plain text content to send.
// - embed: optional embed struct to include.
// - edit: optional message ID to edit.
//
// Returns the sent or edited message object, and any error encountered.
func Respond(e types.Event, content string, embed *types.Embed, edit *string) (interface{}, error) {
	switch e.Platform {
	case "Discord":
		session := e.Session.(*discordgo.Session)
		switch ctx := e.Context.(type) {
		case *discordgo.MessageCreate:
			if edit != nil {
				editMsg := &discordgo.MessageEdit{
					ID:      *edit,
					Channel: ctx.ChannelID,
					Content: &content,
				}
				if embed != nil {
					embeds := []*discordgo.MessageEmbed{{Title: embed.Title, Description: embed.Description, Color: embed.Color}}
					editMsg.Embeds = &embeds
				}
				return session.ChannelMessageEditComplex(editMsg)
			}
			msg := &discordgo.MessageSend{Content: content, Reference: &discordgo.MessageReference{MessageID: ctx.ID, ChannelID: ctx.ChannelID, GuildID: ctx.GuildID}}
			if embed != nil {
				msg.Embeds = []*discordgo.MessageEmbed{{Title: embed.Title, Description: embed.Description, Color: embed.Color}}
			}
			return session.ChannelMessageSendComplex(ctx.ChannelID, msg)

		case *discordgo.InteractionCreate:
			if edit != nil {
				editData := &discordgo.WebhookEdit{Content: &content}
				if embed != nil {
					embeds := []*discordgo.MessageEmbed{{Title: embed.Title, Description: embed.Description, Color: embed.Color}}
					editData.Embeds = &embeds
				}
				return session.InteractionResponseEdit(ctx.Interaction, editData)
			}
			data := &discordgo.InteractionResponseData{Content: content}
			if embed != nil {
				data.Embeds = []*discordgo.MessageEmbed{{Title: embed.Title, Description: embed.Description, Color: embed.Color}}
			}
			return nil, session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: data,
			})
		}

	case "Revolt":
		session := e.Session.(*revoltgo.Session)
		var channelID string
		switch ctx := e.Context.(type) {
		case *revoltgo.EventMessage:
			channelID = ctx.Channel
		case *revoltgo.EventMessageUpdate:
			channelID = ctx.Data.Channel
		case *revoltgo.EventMessageReact, *revoltgo.EventMessageRemoveReaction:
			channelID = ctx.(interface{ GetChannelID() string }).GetChannelID()
		default:
			return nil, fmt.Errorf("unsupported Revolt context type")
		}
		if edit != nil {
			editMsg := revoltgo.MessageEditData{Content: content}
			if embed != nil {
				editMsg.Embeds = []*revoltgo.MessageEmbed{{Title: embed.Title, Description: embed.Description, Colour: fmt.Sprintf("#%06X", embed.Color)}}
			}
			return session.ChannelMessageEdit(channelID, *edit, editMsg)
		}
		sendMsg := revoltgo.MessageSend{Content: content}
		if embed != nil {
			sendMsg.Embeds = []*revoltgo.MessageEmbed{{Title: embed.Title, Description: embed.Description, Colour: fmt.Sprintf("#%06X", embed.Color)}}
		}
		return session.ChannelMessageSend(channelID, sendMsg)
	}
	return nil, fmt.Errorf("unsupported platform or context")
}

// SendMessage sends a basic message to a given channel on either platform.
func SendMessage(platform string, channelID string, content string, embed *types.Embed) (interface{}, error) {
	switch platform {
	case "Discord":
		msg := &discordgo.MessageSend{Content: content}
		if embed != nil {
			msg.Embeds = []*discordgo.MessageEmbed{{Title: embed.Title, Description: embed.Description, Color: embed.Color}}
		}
		return Discord.ChannelMessageSendComplex(channelID, msg)

	case "Revolt":
		sendMsg := revoltgo.MessageSend{Content: content}
		if embed != nil {
			sendMsg.Embeds = []*revoltgo.MessageEmbed{{Title: embed.Title, Description: embed.Description, Colour: fmt.Sprintf("#%06X", embed.Color)}}
		}
		return Revolt.ChannelMessageSend(channelID, sendMsg)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}

// RegisterSlashCommands registers slash commands globally or per guild.
func RegisterSlashCommands(session *discordgo.Session, appID, guildID string, commands []*discordgo.ApplicationCommand) ([]*discordgo.ApplicationCommand, error) {
	if guildID != "" {
		return session.ApplicationCommandBulkOverwrite(appID, guildID, commands)
	}
	return session.ApplicationCommandBulkOverwrite(appID, "", commands)
}

// DeleteMessage deletes a message from a specified platform and channel.
func DeleteMessage(platform string, channelID, messageID string) error {
	switch platform {
	case "Discord":
		return Discord.ChannelMessageDelete(channelID, messageID)
	case "Revolt":
		return Revolt.ChannelMessageDelete(channelID, messageID)
	default:
		return fmt.Errorf("unsupported platform: %s", platform)
	}
}

// EditMessage edits an existing message on either platform.
func EditMessage(platform string, channelID, messageID, content string, embed *types.Embed) (interface{}, error) {
	if platform == "Discord" {
		edit := &discordgo.MessageEdit{ID: messageID, Channel: channelID, Content: &content}
		if embed != nil {
			embeds := []*discordgo.MessageEmbed{{Title: embed.Title, Description: embed.Description, Color: embed.Color}}
			edit.Embeds = &embeds
		}
		return Discord.ChannelMessageEditComplex(edit)
	} else if platform == "Revolt" {
		edit := revoltgo.MessageEditData{Content: content}
		if embed != nil {
			edit.Embeds = []*revoltgo.MessageEmbed{{Title: embed.Title, Description: embed.Description, Colour: fmt.Sprintf("#%06X", embed.Color)}}
		}
		return Revolt.ChannelMessageEdit(channelID, messageID, edit)
	}
	return nil, fmt.Errorf("unsupported platform")
}

// DeferInteraction defers a Discord interaction response, typically used for long processing.
func DeferInteraction(session *discordgo.Session, interaction *discordgo.Interaction) error {
	return session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

// GetUserAvatarURL returns the avatar URL of a user.
func GetUserAvatarURL(user types.User) string {
	return user.Avatar
}

// ParseCommand splits a string into command and arguments.
func ParseCommand(input string) (string, []string) {
	split := strings.Fields(input)
	if len(split) == 0 {
		return "", nil
	}
	return split[0], split[1:]
}

// LogEvent logs an event using the standard logger.
func LogEvent(e types.Event) {
	log.Printf("[%s:%s] %+v\n", e.Platform, e.Type, e)
}

// IsBotEvent returns true if the event came from a bot.
func IsBotEvent(e types.Event) bool {
	return e.Bot
}

// GetAuthor extracts the user object from the event.
func GetAuthor(e types.Event) types.User {
	switch data := e.Data.(type) {
	case types.MessageCallback:
		return data.Author
	case types.InteractionCallback:
		return data.Author
	default:
		return types.User{}
	}
}

// GetChannelID returns the channel ID from an event context.
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
	default:
		return ""
	}
}

// NewCommand creates a new Discord slash command.
func NewCommand(name, desc string, opts ...*discordgo.ApplicationCommandOption) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        name,
		Description: desc,
		Options:     opts,
	}
}

// NewOption creates a new slash command option.
func NewOption(name, desc string, optType discordgo.ApplicationCommandOptionType, required bool) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        name,
		Description: desc,
		Type:        optType,
		Required:    required,
	}
}

// PaginateReply sends multiple messages to a channel with delay in between.
func PaginateReply(session *discordgo.Session, channelID string, pages []string, delay time.Duration) {
	for _, page := range pages {
		_, _ = session.ChannelMessageSend(channelID, page)
		time.Sleep(delay)
	}
}

// FindCommandByName searches a command list for a command with the given name.
func FindCommandByName(commands []*discordgo.ApplicationCommand, name string) *discordgo.ApplicationCommand {
	for _, cmd := range commands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

// EnsureSlashCommands checks for missing commands and registers them if needed.
func EnsureSlashCommands(session *discordgo.Session, appID, guildID string, commands []*discordgo.ApplicationCommand) ([]*discordgo.ApplicationCommand, error) {
	existing, err := session.ApplicationCommands(appID, guildID)
	if err != nil {
		return nil, err
	}

	var toCreate []*discordgo.ApplicationCommand
	for _, cmd := range commands {
		if FindCommandByName(existing, cmd.Name) == nil {
			toCreate = append(toCreate, cmd)
		}
	}

	if len(toCreate) > 0 {
		return session.ApplicationCommandBulkOverwrite(appID, guildID, commands)
	}
	return existing, nil
}

// GetGuildID extracts the guild ID from the event context.
func GetGuildID(e types.Event) string {
	switch ctx := e.Context.(type) {
	case *discordgo.MessageCreate:
		return ctx.GuildID
	case *discordgo.InteractionCreate:
		return ctx.GuildID
	default:
		return ""
	}
}

// GetUsername returns a formatted string with username and ID.
func GetUsername(e types.Event) string {
	user := GetAuthor(e)
	return fmt.Sprintf("%s (%s)", user.Username, user.ID)
}

// IsAdmin checks if the user ID is in the admin list.
func IsAdmin(userID string, adminIDs []string) bool {
	return slices.Contains(adminIDs, userID)
}

// Cooldown returns true if the key is still cooling down. Otherwise, starts a cooldown.
func Cooldown(key string, duration time.Duration) bool {
	cooldownLock.Lock()
	defer cooldownLock.Unlock()

	if until, exists := cooldowns[key]; exists && time.Now().Before(until) {
		return true
	}
	cooldowns[key] = time.Now().Add(duration)
	return false
}

// Retry retries a function for a number of attempts with a 500ms delay.
func Retry(fn func() error, attempts int) error {
	var err error
	for range attempts {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return err
}

// Background runs a function in a new goroutine.
func Background(fn func()) {
	go fn()
}
