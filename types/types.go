package types

// EventType represents the type of event received from a chat platform.
type EventType string

const (
	// MessageCreate is triggered when a new message is created.
	MessageCreate EventType = "MessageCreate"

	// MessageUpdate is triggered when an existing message is edited.
	MessageUpdate EventType = "MessageUpdate"

	// MessageDelete is triggered when a message is deleted.
	MessageDelete EventType = "MessageDelete"

	// ReactionAdd is triggered when a reaction is added to a message.
	ReactionAdd EventType = "ReactionAdd"

	// ReactionRemove is triggered when a reaction is removed from a message.
	ReactionRemove EventType = "ReactionRemove"

	// InteractionCreate is triggered when an interaction is received (e.g., slash command).
	InteractionCreate EventType = "InteractionCreate"
)

// Event holds information about a received platform event.
type Event struct {
	// Name is an optional identifier for the event (custom or descriptive).
	Name string

	// Type specifies the type of event (e.g., MessageCreate).
	Type EventType

	// Platform indicates which platform the event originated from (e.g., "Discord", "Revolt").
	Platform string

	// Bot is true if the event was triggered by the bot itself.
	Bot bool

	// Context is the raw platform-specific context (e.g., *discordgo.MessageCreate or *revoltgo.EventMessage).
	Context any

	// Session is the active session for the platform (e.g., *discordgo.Session).
	Session any

	// Data holds the parsed payload, such as a MessageCallback or InteractionCallback.
	Data interface{}
}

// ActivityType represents the type of activity a bot can display in its presence.
type ActivityType int

const (
	// ActivityTypeGame shows the bot as "Playing [name]".
	ActivityTypeGame ActivityType = iota

	// ActivityTypeStreaming shows the bot as "Streaming [name]".
	ActivityTypeStreaming

	// ActivityTypeListening shows the bot as "Listening to [name]".
	ActivityTypeListening

	// ActivityTypeWatching shows the bot as "Watching [name]".
	ActivityTypeWatching

	// ActivityTypeCustom allows a custom status message.
	ActivityTypeCustom

	// ActivityTypeCompeting shows the bot as "Competing in [name]".
	ActivityTypeCompeting
)

// DiscordConfig stores Discord-specific bot credentials and settings.
type DiscordConfig struct {
	ClientID     string // OAuth2 client ID
	ClientSecret string // OAuth2 client secret
	Token        string // Bot token
}

// RevoltConfig stores Revolt-specific bot credentials.
type RevoltConfig struct {
	Token string // Bot token
}

// Config holds configuration for all supported platforms.
type Config struct {
	Discord DiscordConfig
	Revolt  RevoltConfig
}

// Presence represents a bot's visibility or availability status.
type Presence string

const (
	// Online indicates the bot is active and visible.
	Online Presence = "Online"

	// Idle indicates the bot is inactive but still online.
	Idle Presence = "Idle"

	// DND means "Do Not Disturb"; the bot won't respond.
	DND Presence = "DND"

	// Busy is a custom presence indicating the bot is occupied.
	Busy Presence = "Busy"

	// Invisible appears offline to users.
	Invisible Presence = "Invisible"
)

// User defines a basic user structure returned by chat APIs.
type User struct {
	ID       string // Unique user ID
	Username string // Display name or tag
	Avatar   string // URL or identifier for the user's avatar
}

// MessageCallback represents the content and author of a received message.
type MessageCallback struct {
	Content string // Message text content
	Author  User   // User who sent the message
}

// InteractionCallback represents data sent with a user interaction, such as a slash command.
type InteractionCallback struct {
	Name   string            // Command or interaction name
	Fields map[string]string // Key-value pairs of submitted form data or arguments
	Author User              // User who initiated the interaction
}

// Embed defines a rich content structure for sending styled messages.
type Embed struct {
	Title       string // Title of the embed
	Description string // Body content of the embed
	Color       int    // Color accent (integer representing a hex code)
}
