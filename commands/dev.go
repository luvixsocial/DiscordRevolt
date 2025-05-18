package commands

import (
	"github.com/luvixsocial/whiskercat/types"
)

func EnableDev(evt types.Event, stdout *bool) {
	*stdout = true
	Respond(evt, "Enabled developer mode.", nil, nil)
}

func DisableDev(evt types.Event, stdout *bool) {
	*stdout = false
	Respond(evt, "Disabled developer mode.", nil, nil)
}
