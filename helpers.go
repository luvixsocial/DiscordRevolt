package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Helper function to convert []*discordgo.ApplicationCommandInteractionDataOption to map[string]string
func convertOptionsToMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]string {
	fields := make(map[string]string)
	for _, option := range options {
		if option != nil && option.Value != nil {
			fields[option.Name] = fmt.Sprintf("%v", option.Value)
		}
	}
	return fields
}
