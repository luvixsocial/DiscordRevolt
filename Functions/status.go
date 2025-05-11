package functions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"

	"whiskercat/types"
)

// Set the status of both Discord and Revolt clients.
func SetStatus(ActivityType types.ActivityType, ActivityName string, Presence types.Presence, Status *string) {
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
			statusData.Status = *Status
		}

		err := Discord.UpdateStatusComplex(statusData)
		if err != nil {
			fmt.Println("Failed to set Discord status:", err)
		}
	}

	if Revolt != nil {
		self := Revolt.State.Self()
		if self == nil {
			fmt.Println("Failed to set Revolt status: Self() returned nil")
			return
		}

		status := &revoltgo.UserStatus{
			Text:     ActivityName,
			Presence: revoltgo.UserStatusPresence(Presence),
		}

		_, err := Revolt.UserEdit(self.ID, revoltgo.UserEditData{
			Status: status,
		})
		if err != nil {
			fmt.Println("Failed to set Revolt status:", err)
		}
	}
}
