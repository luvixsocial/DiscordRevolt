package commands

func Test(evt Event, _ *bool) {
	Respond(evt, "Received test event!", nil, nil)
}
