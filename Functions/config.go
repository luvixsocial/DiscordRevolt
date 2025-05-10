package functions

import (
	"libdozina/types"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

// Configure both Discord and Revolt clients.
func Config(config *types.Config) {
	var err error
	Discord, err = discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		panic(err)
	}

	Revolt = revoltgo.New(config.Revolt.Token)
}
