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
	server, f1 := RESTServer.NewRServer(conf)

	//defer the close till the shell has closed
	defer f1.Close()

	if ConfigFilePath != nil && *ConfigFilePath != "" {
		// load this first
		server.ConfigFilePath = *ConfigFilePath
		ok := server.LoadConfig()

		if !ok {
			return
		}
	}

	RESTCommands.AddDefaults(server)
	RESTCommands.SetDefaultFunctionalMap(server)

	// start Listen Server, this build the mapping and creates Handler/
	// also fires the "http listen and server method"
	server.ListenAndServe()

}
