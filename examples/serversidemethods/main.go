package main

import (
	RESTCommands "github.com/eshu0/RESTServer/pkg/commands"
	RESTConfig "github.com/eshu0/RESTServer/pkg/config"
	RESTServer "github.com/eshu0/RESTServer/pkg/server"
)

func main() {

	// create a new server
	conf := RESTConfig.NewRServerConfig()

	// Create a new REST Server
	server := RESTServer.NewRServer(conf)

	//defer the close till the server has closed
	defer server.FinishLogging()

	// load this first
	server.ConfigFilePath = "./config.json"

	ok := server.LoadConfig()

	if !ok {
		return
	}

	// add the defaults here
	RESTCommands.AddDefaults(server)
	RESTCommands.SetDefaultFunctionalMap(server)

	// this registers the custom structures
	// in the JSON config the FunctionalClass is the name used for the map "TestAnother"
	// if these are not public and spelt correctly the lookups will fail
	server.RegisterFunction("TestAnother", TestAnother{})
	server.RegisterFunction("TestStruct", TestStruct{})
	server.Register("TypedTestAnother", TypedTestAnother{})
	server.Register("TypedTestStruct", TypedTestStruct{})

	// as a test save the updated config
	server.ConfigFilePath = "./updated.json"
	server.SaveConfig()

	// start Listen Server, this build the mapping and creates Handler/
	// also fires the "http listen and server method"
	server.ListenAndServe()

}
