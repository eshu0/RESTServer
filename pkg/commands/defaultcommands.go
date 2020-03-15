package RESTCommands

import (
	"errors"
	"html/template"
	"net/http"

	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	Server "github.com/eshu0/RESTServer/pkg/server"
)

type RServerCommand struct {
	Server *Server.RServer
}

func checkRCommand(rsc RServerCommand) bool {
	if rsc.Server != nil {
		return false
	}

	return true
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

		t := template.New("old") // *Template{}
		err := errors.New("should not see this error")

		if rsc.Server.Config.HasTemplate() {
			t, err = template.ParseFiles(rsc.Server.Config.GetTemplatePath())
		} else {
			doc := "<!DOCTYPE html><html><head><meta charset=\"UTF-8\"><title>Available REST API Calls</title></head><body><h1>Available REST API Calls</h1><h2>Custom</h2>{{range .Handlers}}<div><a href=\"{{ .URL }}\">{{ .URL }} will point to {{ .HTTPMethod }} {{ .FunctionalClass }}.{{ .MethodName }} </a></div>{{else}}<div><strong>None</strong></div>{{end}}<h2>Default</h2>{{range .DefaultHandlers}}<div><a href=\"{{ .URL }}\">{{ .MethodName }}</a></div></div>{{else}}<div><strong>No Handlers</strong></div>{{end}}</body></html>"
			t, err = template.New("list").Parse(doc)
		}

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

func (rsc RServerCommand) CreateShutDownHandler() Handlers.RESTHandler {
	drhs := Handlers.RESTHandler{}

	drhs.URL = "/admin/shutdown"
	drhs.MethodName = "ShutDown"
	drhs.HTTPMethod = "GET"
	drhs.FunctionalClass = "RServerCommand"
	return drhs
}

func (rsc RServerCommand) CreateListHandler() Handlers.RESTHandler {
	drhr := Handlers.RESTHandler{}

	drhr.URL = "/admin/listcommands"
	drhr.MethodName = "ListCommands"
	drhr.HTTPMethod = "GET"
	drhr.FunctionalClass = "RServerCommand"
	return drhr
}

func (rsc RServerCommand) CreateLoadConfigHandler() Handlers.RESTHandler {
	drhr := Handlers.RESTHandler{}

	drhr.URL = "/admin/loadconfig"
	drhr.MethodName = "LoadConfig"
	drhr.HTTPMethod = "GET"
	drhr.FunctionalClass = "RServerCommand"
	return drhr
}

func (rsc RServerCommand) CreateSaveConfigHandler() Handlers.RESTHandler {
	drhr := Handlers.RESTHandler{}

	drhr.URL = "/admin/saveconfig"
	drhr.MethodName = "SaveConfig"
	drhr.HTTPMethod = "GET"
	drhr.FunctionalClass = "RServerCommand"
	return drhr
}

func AddDefaults(server *Server.RServer) {

	// Default commands for server
	// These should be removed if not required
	rsc := RServerCommand{Server: server}

	server.FunctionalMap["RServerCommand"] = rsc

	server.Config.AddDefaultHandler(rsc.CreateShutDownHandler())
	server.Config.AddDefaultHandler(rsc.CreateListHandler())
	server.Config.AddDefaultHandler(rsc.CreateLoadConfigHandler())
	server.Config.AddDefaultHandler(rsc.CreateSaveConfigHandler())

	for _, handl := range server.Config.GetDefaultHandlers() {
		server.Log.LogDebugf("NewRServerWithDefaults", "Default Handler: Added %s", handl.MethodName)
	}
}

func SetDefaultFunctionalMap(server *Server.RServer) {
	server.FunctionalMap["RServerCommand"] = RServerCommand{Server: server}
}
