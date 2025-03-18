package functions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

// Configure both Discord and Revolt clients.
func Config(discord string, revolt string) {
	var err error
	Revolt = revoltgo.New(revolt)
	Discord, err = discordgo.New(discord)
	if err != nil {
		panic(err)
	}
}
