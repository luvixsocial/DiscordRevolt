package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/luvixsocial/whiskercat/types"
	"github.com/sentinelb51/revoltgo"
)

// SetStatus updates the bot's activity and presence for both Discord and Revolt.
//
// Parameters:
// - activityType: What the bot is doing (e.g., Playing, Watching)
// - activityName: Name of the activity
// - presence: Presence status (Online, Idle, DND, Invisible)
// - rawDiscordStatus: Optional Discord raw status ("online", "idle", "dnd", "invisible")
func SetStatus(activityType types.ActivityType, activityName string, presence types.Presence, rawDiscordStatus *string) {
	setDiscordStatus(activityType, activityName, rawDiscordStatus)
	setRevoltStatus(activityName, presence)
}

func setDiscordStatus(activityType types.ActivityType, activityName string, status *string) {
	if Discord == nil {
		return
	}

	update := discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: activityName,
				Type: discordgo.ActivityType(activityType),
			},
		},
	}

	if status != nil {
		update.Status = *status
	}

	if err := Discord.UpdateStatusComplex(update); err != nil {
		log.Println("❌ Discord status update failed:", err)
	} else {
		log.Println("✅ Discord status updated.")
	}
}

func setRevoltStatus(activityName string, presence types.Presence) {
	if Revolt == nil {
		return
	}

	self := Revolt.State.Self()
	if self == nil {
		log.Println("❌ Revolt status update failed: Self() is nil")
		return
	}

	status := &revoltgo.UserStatus{
		Text:     activityName,
		Presence: revoltgo.UserStatusPresence(presence),
	}

	if _, err := Revolt.UserEdit(self.ID, revoltgo.UserEditData{
		Status: status,
	}); err != nil {
		log.Println("❌ Revolt status update failed:", err)
	} else {
		log.Println("✅ Revolt status updated.")
	}
}
