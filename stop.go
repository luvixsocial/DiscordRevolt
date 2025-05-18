package main

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
