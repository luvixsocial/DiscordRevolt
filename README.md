# WhiskerCat

WhiskerCat is a powerful Go library designed to handle WebSocket events from both **Discord** and **Revolt Chat**. With **WhiskerCat**, developers can build cross-platform bots that use a single command structure to interact seamlessly with both platforms.

## Features

- **Unified Event Handling** - Receive and process WebSocket events from **Discord** and **Revolt Chat** in one place.
- **Cross-Platform Commands** - Write a single command file that responds to both platforms without extra logic.
- **WebSocket Abstraction** - Simplifies the complexities of WebSocket connections.
- **Lightweight & Efficient** - Optimized for performance and scalability.
- **Easy-to-Use API** - Minimal setup required to start building bots.

## Installation

To install **WhiskerCat**, run the following command:

```sh
go get github.com/luvixsocial/WhiskerCat
```

## Getting Started

### Basic Example

Below is a simple example demonstrating how to initialize the bot and handle events from both platforms:

```go
package main

import (
	"fmt"
	bot "github.com/luvixsocial/WhiskerCat"
)

// Test
func main() {
	bot.Config(&types.Config{
        Discord: &types.DiscordConfig{
            ClientID:     "YOUR_DISCORD_CLIENT_ID",
            ClientSecret: "YOUR_DISCORD_CLIENT_SECRET",
            Token:        "YOUR_DISCORD_BOT_TOKEN",
        },
        Revolt: &types.RevoltConfig{
            Token: "YOUR_REVOLT_BOT_TOKEN",
        },
    })
	bot.Start()
	bot.SetStatus(types.ActivityTypeGame, "Working on Luvix Social", types.Online, nil)

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
 bot.Respond(evt, "Hello from WhiskerCat!", nil, nil)
})
```

### Supported Events
WhiskerCat currently supports handling the following events:

- **Message Create** - Triggered when a message is sent in a channel.
- **Message Update** - Triggered when a message is updated.
- **Message Delete* - Triggered when a message is deleted.
- **ReactionAdd** - Triggered when a message receives a reaction.
- **ReactionRemove** - Triggered when a message loses a reaction.
- **Interaction Create (Discord only)** - Triggered when an interaction (such as a slash command) is executed.

## Status Management

You can set the bot's status using the `SetStatus()` function:

```go
bot.SetStatus(types.ActivityTypeGame, "Developing with WhiskerCat", types.Online, nil)
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

_WhiskerCat is part of the **Luvix Social** ecosystem, providing powerful tools for seamless bot development._

