package types

// Events
type EventType string

const (
	Message     EventType = "Message"
	Interaction EventType = "Interaction"
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

// Activities
type ActivityType int

const (
	ActivityTypeGame      ActivityType = 0
	ActivityTypeStreaming ActivityType = 1
	ActivityTypeListening ActivityType = 2
	ActivityTypeWatching  ActivityType = 3
	ActivityTypeCustom    ActivityType = 4
	ActivityTypeCompeting ActivityType = 5
)

// Presences
type Presence string

const (
	Online    Presence = "Online"
	Idle      Presence = "Idle"
	DND       Presence = "DND"
	Busy      Presence = "Busy"
	Invisible Presence = "Invisible"
)

// Users
type User struct {
	ID       string
	Username string
	Avatar   string
}

// Message Callback
type MessageCallback struct {
	Content string
	Author  User
}

// Interaction Callback
type InteractionCallback struct {
	Name   string
	Fields map[string]string
	Author User
}
