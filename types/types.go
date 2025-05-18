package types

import "github.com/bwmarrin/discordgo"

// EventType defines supported platform event types.
type EventType string

const (
	MessageCreate          EventType = "MessageCreate"
	MessageUpdate          EventType = "MessageUpdate"
	MessageDelete          EventType = "MessageDelete"
	ReactionAdd            EventType = "ReactionAdd"
	ReactionRemove         EventType = "ReactionRemove"
	InteractionCreate      EventType = "InteractionCreate"
	EventTypingStart       EventType = "TypingStart"
	EventVoiceStateUpdate  EventType = "VoiceStateUpdate"
	EventPresenceUpdate    EventType = "PresenceUpdate"
	EventGuildMemberAdd    EventType = "GuildMemberAdd"
	EventGuildMemberRemove EventType = "GuildMemberRemove"
	EventChannelCreate     EventType = "ChannelCreate"
	EventChannelUpdate     EventType = "ChannelUpdate"
	EventChannelDelete     EventType = "ChannelDelete"
	EventUserUpdate        EventType = "UserUpdate"
	EventMemberJoin        EventType = "MemberJoin"  // Revolt specific
	EventMemberLeave       EventType = "MemberLeave" // Revolt specific
)

// Event represents a normalized platform event.
type Event struct {
	Name     string    // Optional identifier for the event
	Type     EventType // The type of event triggered
	Platform string    // "Discord" or "Revolt"
	Bot      bool      // True if the event was triggered by a bot
	Context  any       // The raw platform event (e.g., *discordgo.MessageCreate)
	Session  any       // The session for the platform
	Data     any       // Parsed payload like MessageCallback or InteractionCallback
}

// ActivityType describes what the bot is shown doing.
type ActivityType int

const (
	ActivityTypeGame      ActivityType = iota // Playing
	ActivityTypeStreaming                     // Streaming
	ActivityTypeListening                     // Listening to
	ActivityTypeWatching                      // Watching
	ActivityTypeCustom                        // Custom status
	ActivityTypeCompeting                     // Competing in
)

// Presence defines bot's online visibility.
type Presence string

const (
	Online    Presence = "Online"
	Idle      Presence = "Idle"
	DND       Presence = "DND"
	Busy      Presence = "Busy"
	Invisible Presence = "Invisible"
)

// DiscordConfig contains Discord credentials.
type DiscordConfig struct {
	ClientID     string // OAuth2 Client ID
	ClientSecret string // OAuth2 Client Secret
	Token        string // Bot Token
}

// RevoltConfig contains Revolt credentials.
type RevoltConfig struct {
	Token string // Bot Token
}

// AuthConfig aggregates credentials for all platforms.
type AuthConfig struct {
	Discord DiscordConfig // Discord configuration
	Revolt  RevoltConfig  // Revolt configuration
}

// User represents a basic user identity.
type User struct {
	ID       string // Unique user ID
	Username string // Display name
	Avatar   string // Avatar URL or identifier
}

// MessageCallback wraps a received message and author.
type MessageCallback struct {
	Content string // Message content
	Author  User   // Message author
}

// InteractionCallback holds interaction data like slash commands.
type InteractionCallback struct {
	Name   string                       // Name of the interaction/command
	Fields map[string]string            // Option key-value map
	Data   *discordgo.InteractionCreate // Raw interaction object (Discord only)
	Author User                         // Command invoker
}

// Embed defines a structured rich message.
type Embed struct {
	Title       string        // Embed title
	Description string        // Main content body
	IconURL     *string       // Icon URL
	PhotoURL    *string       // Image URL
	URL         *string       // Link to the embed
	Fields      *[]EmbedField // List of fields
	Footer      *EmbedFooter  // Footer text
	Color       int           // Accent color as integer (hex)
}
type EmbedFooter struct {
	Text     string // Footer text
	PhotoURL string // Footer Photo URL
}
type EmbedField struct {
	Name   string // Field name
	Value  string // Field value
	Inline bool   // Whether to display inline
}
