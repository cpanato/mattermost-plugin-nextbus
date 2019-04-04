package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/cpanato/mattermost-plugin-nextbus/server/nextbus"
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	botUserID string

	nextBusClient *nextbus.Client
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
