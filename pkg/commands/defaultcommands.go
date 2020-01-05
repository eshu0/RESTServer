package RESTServer

import (
	"html/template"
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
		rsc.Server.Log.LogDebug("RServerCommand", "List Commands called")

		doc := "<!DOCTYPE html><html><head><meta charset=\"UTF-8\"><title>Available REST API Calls</title></head><body><h1>Available REST API Calls</h1><h2>Custom</h2>{{range .Handlers}}<div><a href=\"{{ .URL }}\">{{ .MethodName }}</a></div>{{else}}<div><strong>None</strong></div>{{end}}<h2>Default</h2>{{range .DefaultHandlers}}<div><a href=\"{{ .URL }}\">{{ .MethodName }}</a></div></div>{{else}}<div><strong>No Handlers</strong></div>{{end}}</body></html>"
		t, err := template.New("list").Parse(doc)
		if err != nil {
			rsc.Server.Log.LogErrorf("RServerCommand", "Error : %s", err.Error())
			return
		}
		err = t.Execute(w, rsc.Server.Config) // Template(w, "T", "<script>alert('you have been pwned')</script>")
		if err != nil {
			rsc.Server.Log.LogErrorf("RServerCommand", "Error : %s", err.Error())
			return
		}
	}
}

func (rsc RServerCommand) LoadConfig(w http.ResponseWriter, r *http.Request) {
	if rsc.Server != nil {
		rsc.Server.Log.LogDebug("RServerCommand", "Load Config called")
		rsc.Server.LoadConfig()
	}
}

func (rsc RServerCommand) SaveConfig(w http.ResponseWriter, r *http.Request) {
	if rsc.Server != nil {
		rsc.Server.Log.LogDebug("RServerCommand", "Save Config called")
		rsc.Server.SaveConfig()
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
