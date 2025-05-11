package types

/*
 * @title Events
 * @description These types are used to define the events that can be received from the Discord and Revolt APIs.
 */
type EventType string

const (
	Message       EventType = "Message"
	MessageUpdate EventType = "MessageUpdate"
	Interaction   EventType = "Interaction"
)

type Event struct {
	Name     string
	Type     EventType
	Platform string
	Bot      bool
	Context  any
	Session  any
	Data     interface{}
}

/*
 * @title Activity Types
 * @description These types are used to define the activity types that can be set for the bots presence.
 */
type ActivityType int

const (
	ActivityTypeGame      ActivityType = 0
	ActivityTypeStreaming ActivityType = 1
	ActivityTypeListening ActivityType = 2
	ActivityTypeWatching  ActivityType = 3
	ActivityTypeCustom    ActivityType = 4
	ActivityTypeCompeting ActivityType = 5
)

/*
 * @title Config
 * @description These types are used to define the configuration for the bots.
 */
type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	Token        string
}

type RevoltConfig struct {
	Token string
}

type Config struct {
	Discord DiscordConfig
	Revolt  RevoltConfig
}

/*
 * @title Presence
 * @description These types are used to define the presence of the bots.
 */
type Presence string

const (
	Online    Presence = "Online"
	Idle      Presence = "Idle"
	DND       Presence = "DND"
	Busy      Presence = "Busy"
	Invisible Presence = "Invisible"
)

/*
 * @title User
 * @description These types are used to define the user object that is returned from the Discord and Revolt APIs.
 */
type User struct {
	ID       string
	Username string
	Avatar   string
}

/*
 * @title Message
 * @description These types are used to define the message object that is returned from the Discord and Revolt APIs.
 */
type MessageCallback struct {
	Content string
	Author  User
}

/*
 * @title Interaction
 * @description These types are used to define the interaction object that is returned from the Discord and Revolt APIs.
 */
type InteractionCallback struct {
	Name   string
	Fields map[string]string
	Author User
}

/*
 * @title Embed
 * @description These types are used to define the embed object that is sent to the Discord and Revolt APIs.
 */
type Embed struct {
	Title       string
	Description string
	Color       int
}
