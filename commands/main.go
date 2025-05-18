package commands

import (
	"fmt"
	"time"

	"github.com/luvixsocial/whiskercat/types"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

var commandMap = map[string]func(types.Event, *bool){
	"ping":        Ping,
	"test":        Test,
	"test_embed":  TestEmbed,
	"enable_dev":  EnableDev,
	"disable_dev": DisableDev,
}

func Handle(evt types.Event, stdout *bool, name string) {
	if handler, ok := commandMap[name]; ok {
		handler(evt, stdout)
	} else if *stdout {
		Respond(evt, "", &Embed{
			Title:       "types.Event Received",
			Description: fmt.Sprintf("%+v", evt),
			Color:       0x00FF00,
		}, nil)
	}
}

func Definitions() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Check latency",
		},
		{
			Name:        "test",
			Description: "Test types.event reply",
		},
		{
			Name:        "test_embed",
			Description: "Send a test embed",
		},
		{
			Name:        "enable_dev",
			Description: "Enable developer mode",
		},
		{
			Name:        "disable_dev",
			Description: "Disable developer mode",
		},
	}
}
