package main

import (
	"github.com/eshu0/RESTServer/pkg/server"
	"github.com/eshu0/RESTServer/pkg/config"
	"github.com/eshu0/RESTServer/pkg/handlers"
)

func main() {

	conf := RESTServer.NewRServerConfig()

	// Create a new REST Server
	server, f1 := RESTServer.NewRServer(conf)

	//defer the close till the shell has closed
	defer f1.Close()

	// load this first
	server.ConfigFilePath = "./config.json"
	ok := server.LoadConfig()

	if !ok {
		return
	}

	server.AddDefaults()

	// save the updated config
	server.ConfigFilePath = "./updated.config"
	server.SaveConfig()

	// start Listen Server, this build the mapping and creates Handler/
	// also fires the "http listen and server method"
	server.ListenAndServe()

}
