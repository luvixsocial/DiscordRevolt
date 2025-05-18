package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/luvixsocial/whiskercat/types"
	"github.com/sentinelb51/revoltgo"
	"log"
)

// Config sets up and initializes both Discord and Revolt clients using the provided configuration.
//
// This should be called before any event handlers or client operations.
// It panics if Discord client initialization fails.
func Config(config *types.AuthConfig) {
	var err error

	// Initialize Discord client
	Discord, err = discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		log.Fatalf("Failed to initialize Discord client: %v", err)
	}

	// Initialize Revolt client
	Revolt = revoltgo.New(config.Revolt.Token)
}
