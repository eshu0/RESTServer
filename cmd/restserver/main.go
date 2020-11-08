package main

import (
	"flag"

	RESTServer "github.com/eshu0/RESTServer/pkg/server"
)

func main() {

	ConfigFilePath := flag.String("config", "", "Filepath to config file")
	UpdatedConfigFilePath := flag.String("newconfigpath", "", "Filepath to config file to be saved after ")
	flag.Parse()

	logger := sl.NewApplicationLogger()

	server := RESTServer.DefaultServer(ConfigFilePath, logger)

	// add the defaults here
	RESTCommands.AddDefaults(server)
	RESTCommands.SetDefaultFunctionalMap(server)

	// start Listen Server, this build the mapping and creates Handler/
	// also fires the "http listen and server method"
	server.ListenAndServe()

	// has a conifg file been provided?
	if UpdatedConfigFilePath != nil && *UpdatedConfigFilePath != "" {

		// as a test save the updated config
		server.ConfigFilePath = UpdatedConfigFilePath
		server.SaveConfig()
	}

}
