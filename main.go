package main

// Importing the required packages
import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

// Declaring the global variables
var (
	Revolt  *revoltgo.Session
	Discord *discordgo.Session
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

// Configure both Discord and Revolt clients.
func Config(discord string, revolt string) {
	var err error
	Revolt = revoltgo.New(revolt)
	Discord, err = discordgo.New(discord)
	if err != nil {
		panic(err)
	}
}

// Start both Discord and Revolt clients.
func Start() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		err := Discord.Open()
		if err != nil {
			log.Fatalf("Error starting Discord client: %v", err)
		}
		fmt.Println("Discord client started!")
		defer wg.Done()
	}()

	go func() {
		err := Revolt.Open()
		if err != nil {
			log.Fatalf("Error starting Revolt client: %v", err)
		}
		fmt.Println("Revolt client started!")
		time.Sleep(1 * time.Second)
		defer wg.Done()
	}()

	wg.Wait()
}

// Stop both Discord and Revolt clients.
func Stop() {
	err := Discord.Close()
	if err != nil {
		panic("Error stopping Discord client: " + err.Error())
	}

	err = Revolt.Close()
	if err != nil {
		panic("Error stopping Revolt client: " + err.Error())
	}
}

// Set the status of both Discord and Revolt clients.
func SetStatus(ActivityType ActivityType, ActivityName string, Presence Presence, Status *string) {
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

// Event Listener for both Discord and Revolt clients.
func OnEvent(callback func(Event)) {
	if Discord != nil {
		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
			eventType := fmt.Sprintf("%T", e)
			callback(Event{
				Name: eventType,
				Type: Message,
				Data: MessageCallback{
					Content: e.Message.Content,
					Author: User{
						ID:       e.Author.ID,
						Username: e.Author.Username,
						Avatar:   e.Author.AvatarURL("128"),
					},
				},
			})
		})

		Discord.AddHandler(func(s *discordgo.Session, e *discordgo.InteractionCreate) {
			eventType := fmt.Sprintf("%T", e)
			callback(Event{
				Name:     eventType,
				Type:     Interaction,
				Platform: "Discord",
				Data: InteractionCallback{
					Name:   e.ApplicationCommandData().Name,
					Fields: convertOptionsToMap(e.ApplicationCommandData().Options),
					Author: User{
						ID:       e.Member.User.ID,
						Username: e.Member.User.Username,
						Avatar:   e.Member.User.AvatarURL("128"),
					},
				},
			})
		})
	}

	if Revolt != nil {
		Revolt.AddHandler(func(e *revoltgo.Session, m *revoltgo.EventMessage) {
			eventType := fmt.Sprintf("%T", m)
			authorData, err := Revolt.User(m.Author)

			if err != nil {
				fmt.Println("Failed to get author data:", err)
				return
			}

			callback(Event{
				Name:     eventType,
				Type:     Message,
				Platform: "Revolt",
				Data: MessageCallback{
					Content: m.Content,
					Author: User{
						ID:       authorData.ID,
						Username: authorData.Username,
						Avatar:   authorData.Avatar.URL("128"),
					},
				},
			})
		})
	}
}
