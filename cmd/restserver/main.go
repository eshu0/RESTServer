package main

import (
	"flag"

	"github.com/eshu0/RESTServer/pkg/commands"
	"github.com/eshu0/RESTServer/pkg/config"
	"github.com/eshu0/RESTServer/pkg/server"
)

func main() {

	ConfigFilePath := flag.String("config", "", "Filepath to config file")
	flag.Parse()

	conf := RESTConfig.NewRServerConfig()

	// Create a new REST Server
	server := RESTServer.NewRServer(conf)

	log := server.Log
	//defer the close till the shell has closed
	defer log.CloseAllChannels()

	if ConfigFilePath != nil && *ConfigFilePath != "" {
		// load this first
		server.ConfigFilePath = *ConfigFilePath
		ok := server.LoadConfig()

		if !ok {
			return
		}
	} else {
		// load this first
		server.ConfigFilePath = "./config.json"
	}

	RESTCommands.AddDefaults(server)
	RESTCommands.SetDefaultFunctionalMap(server)
	server.SaveConfig()

	// start Listen Server, this build the mapping and creates Handler/
	// also fires the "http listen and server method"
	server.ListenAndServe()

}
