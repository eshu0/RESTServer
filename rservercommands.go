package RESTServer

import (
	"net/http"
)

type RServerCommand struct {
	Server *RServer
}

func (rsc RServerCommand) ShutDown(w http.ResponseWriter, r *http.Request) {
	if rsc.Server != nil {
		rsc.Server.Log.LogDebug("RServerCommand", "HTTP server shutdown called")
		rsc.Server.ShutDown()
	}

}

func (rsc RServerCommand) ListCommands(w http.ResponseWriter, r *http.Request) {
	if rsc.Server != nil {
		rsc.Server.Log.LogDebug("RServerCommand", "HTTP server restart called")

	}
}

func (rsc RServerCommand) LoadConfig(w http.ResponseWriter, r *http.Request) {
	if rsc.Server != nil {
		rsc.Server.Log.LogDebug("RServerCommand", "Load Config called")
		rsc.Server.LoadJSONFile()
	}
}

func (rsc RServerCommand) SaveConfig(w http.ResponseWriter, r *http.Request) {
	if rsc.Server != nil {
		rsc.Server.Log.LogDebug("RServerCommand", "Save Config called")
		rsc.Server.SaveJSONFile()
	}
}

func (rsc RServerCommand) CreateShutDownHandler() RESTHandler {
	drhs := RESTHandler{}

	drhs.URL = "/admin/shutdown"
	drhs.MethodName = "ShutDown"
	drhs.HTTPMethod = "GET"
	drhs.FunctionalClass = "RServerCommand"
	return drhs
}

func (rsc RServerCommand) CreateListHandler() RESTHandler {
	drhr := RESTHandler{}

	drhr.URL = "/admin/listcommands"
	drhr.MethodName = "ListCommands"
	drhr.HTTPMethod = "GET"
	drhr.FunctionalClass = "RServerCommand"
	return drhr
}

func (rsc RServerCommand) CreateLoadConfigHandler() RESTHandler {
	drhr := RESTHandler{}

	drhr.URL = "/admin/loadconfig"
	drhr.MethodName = "LoadConfig"
	drhr.HTTPMethod = "GET"
	drhr.FunctionalClass = "RServerCommand"
	return drhr
}

func (rsc RServerCommand) CreateSaveConfigHandler() RESTHandler {
	drhr := RESTHandler{}

	drhr.URL = "/admin/saveconfig"
	drhr.MethodName = "SaveConfig"
	drhr.HTTPMethod = "GET"
	drhr.FunctionalClass = "RServerCommand"
	return drhr
}
