# DiscordRevolt

DiscordRevolt is a powerful Go library designed to handle WebSocket events from both **Discord** and **Revolt Chat**. With **DiscordRevolt**, developers can build cross-platform bots that use a single command structure to interact seamlessly with both platforms.

## Features

- **Unified Event Handling** - Receive and process WebSocket events from **Discord** and **Revolt Chat** in one place.
- **Cross-Platform Commands** - Write a single command file that responds to both platforms without extra logic.
- **WebSocket Abstraction** - Simplifies the complexities of WebSocket connections.
- **Lightweight & Efficient** - Optimized for performance and scalability.
- **Easy-to-Use API** - Minimal setup required to start building bots.

## Installation

To install **DiscordRevolt**, run the following command:

```sh
go get github.com/luvixsocial/discordrevolt
```

## Getting Started

### Basic Example

Below is a simple example demonstrating how to initialize the bot and handle events from both platforms:

```go
package main

import (
	"fmt"
	bot "github.com/luvixsocial/discordrevolt"
)

// Test
func main() {
	bot.Config("YOUR_DISCORD_TOKEN", "YOUR_REVOLT_TOKEN")
	bot.Start()
	bot.SetStatus(bot.ActivityTypeGame, "Working on Luvix Social", bot.Online, nil)

	bot.OnEvent(func(evt Event) {
		fmt.Printf("Received event: %s\nType: %+v\nData: %+v\n", evt.Name, evt.Type, evt.Data)
	})

	select {}
}
```

## Event Handling

Events from both **Discord** and **Revolt Chat** can be handled using `OnEvent()`.

```go
OnEvent(func(evt Event) {
	fmt.Println("Received event:", evt.Name, "Type:", evt.Type, "Data:", evt.Data)
})
```

### Supported Events
DiscordRevolt currently supports handling the following events:

- **Message Create** - Triggered when a message is sent in a channel.
- **Interaction Create (Discord only)** - Triggered when an interaction (such as a slash command) is executed.

## Status Management

You can set the bot's status using the `SetStatus()` function:

```go
bot.SetStatus(bot.ActivityTypeGame, "Developing with DiscordRevolt", bot.Online, nil)
```

## Advanced Usage

### Custom Event Processing

If you need to handle specific event types differently, you can structure event handling like this:

```go
bot.OnEvent(func(evt Event) {
	switch evt.Name {
	case "MESSAGE_CREATE":
		fmt.Println("New message received:", evt.Data)
	case "INTERACTION_CREATE":
		fmt.Println("New interaction received (Discord only):", evt.Data)
	default:
		fmt.Println("Unhandled event:", evt.Name)
	}
})
```

### Sending Messages

To send a message to a Discord or Revolt channel:

```go
bot.SendMessage("channelID", "Hello, world!")
```

## Contributing

Contributions are welcome! Feel free to submit a pull request or open an issue.

## License

This project is licensed under the **MIT License**.

---

_DiscordRevolt is part of the **Luvix Social** ecosystem, providing powerful tools for seamless bot development._

