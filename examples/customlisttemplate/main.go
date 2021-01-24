package main

import (
	"flag"

	RESTCommands "github.com/eshu0/RESTServer/pkg/commands"
	RESTServer "github.com/eshu0/RESTServer/pkg/server"
)

func main() {

	ConfigFilePath := flag.String("config", "", "Filepath to config file")

	// create a new server - don't pass in a load path
	server := RESTServer.DefaultServer(ConfigFilePath)

	// add the defaults here
	RESTCommands.AddDefaults(server)
	RESTCommands.SetDefaultFunctionalMap(server)

	// as a test save the updated config
	server.Config.Helper.FilePath = "./updated.json"
	server.SaveConfig()

	// start Listen Server, this build the mapping and creates Handler/
	// also fires the "http listen and server method"
	server.ListenAndServe()

}
