package functions

import (
	"fmt"
	"log"
	"sync"
	"time"
)

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
