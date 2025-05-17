package functions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"

	"github.com/luvixsocial/whiskercat/types"
)

// SetStatus updates the activity and presence status for both Discord and Revolt clients.
//
// Parameters:
// - ActivityType: What the bot is doing (e.g., Playing, Watching, Listening)
// - ActivityName: Name of the activity (e.g., "Whisker Adventures")
// - Presence: Presence state (e.g., Online, Idle, DND, Invisible)
// - Status: Optional raw status string for Discord ("online", "idle", "dnd", "invisible")
func SetStatus(ActivityType types.ActivityType, ActivityName string, Presence types.Presence, Status *string) {
	// Discord status update
	if Discord != nil {
		statusData := discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: ActivityName,
					Type: discordgo.ActivityType(ActivityType),
				},
			},
		}

		if Status != nil {
			statusData.Status = *Status // Accepts: "online", "idle", "dnd", "invisible"
		}

		if err := Discord.UpdateStatusComplex(statusData); err != nil {
			fmt.Println("❌ Failed to set Discord status:", err)
		} else {
			fmt.Println("✅ Discord status updated.")
		}
	}

	// Revolt status update
	if Revolt != nil {
		self := Revolt.State.Self()
		if self == nil {
			fmt.Println("❌ Failed to set Revolt status: Self() returned nil")
			return
		}

		status := &revoltgo.UserStatus{
			Text:     ActivityName,
			Presence: revoltgo.UserStatusPresence(Presence), // Accepts: "Online", "Idle", "DND", etc.
		}

		if _, err := Revolt.UserEdit(self.ID, revoltgo.UserEditData{
			Status: status,
		}); err != nil {
			fmt.Println("❌ Failed to set Revolt status:", err)
		} else {
			fmt.Println("✅ Revolt status updated.")
		}
	}
}
