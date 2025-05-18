package commands

import (
	"fmt"
	"time"

	"github.com/luvixsocial/whiskercat/types"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

func Ping(evt types.Event, _ *bool) {
	start := time.Now()
	msgRaw, err := Respond(evt, "Pinging...", nil, nil)
	if err != nil {
		fmt.Printf("Error sending ping: %v\n", err)
		return
	}
	latency := time.Since(start).Milliseconds()
	pong := fmt.Sprintf("Pong! %dms", latency)
	switch msg := msgRaw.(type) {
	case *discordgo.Message:
		Respond(evt, pong, nil, &msg.ID)
	case *revoltgo.Message:
		Respond(evt, pong, nil, &msg.ID)
	}
}
