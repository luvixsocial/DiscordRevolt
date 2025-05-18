package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sentinelb51/revoltgo"
)

// Start initializes and runs both the Discord and Revolt clients concurrently.
// It will block until both clients are successfully started or fatally fail.
//
// This function uses a WaitGroup to ensure both clients are launched before continuing.
// In the event of an error while opening a client, the program will terminate immediately.
func Start() (*discordgo.Session, *revoltgo.Session) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Start Discord client in a separate goroutine
	go func() {
		defer wg.Done()

		if Discord == nil {
			log.Fatalln("Discord client is not initialized")
		}

		if err := Discord.Open(); err != nil {
			log.Fatalf("Error starting Discord client: %v", err)
		}

		fmt.Println("✅ Discord client started!")
	}()

	// Start Revolt client in a separate goroutine
	go func() {
		defer wg.Done()

		if Revolt == nil {
			log.Fatalln("Revolt client is not initialized")
		}

		if err := Revolt.Open(); err != nil {
			log.Fatalf("Error starting Revolt client: %v", err)
		}

		fmt.Println("✅ Revolt client started!")
		time.Sleep(1 * time.Second) // Optional: delay to prevent startup race conditions
	}()

	wg.Wait()

	return Discord, Revolt
}
