package RESTCommands

import (
	"html/template"
	"net/http"

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

func (rsc RServerCommand) ListCommands(w http.ResponseWriter, r *http.Request,t *template.Template) {
	if rsc.Server != nil {
		err := t.Execute(w, rsc.Server.Config) // Template(w, "T", "<script>alert('you have been pwned')</script>")
		if err != nil {
			rsc.Server.Log.LogErrorf("MakeTemplateHandlerFunction", "Error : %s", err.Error())
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

func AddDefaults(server *Server.RServer) {

	dhlen := server.Config.GetDefaultHandlersLen()
	if  dhlen > 0{
		server.Log.LogDebugf("AddDefaults", "Not adding defaults as there are %d handlers already", dhlen)
		return
	}

	// Default commands for server
	// These should be removed if not required
	rsc := RServerCommand{Server: server}

	server.FunctionalMap["RServerCommand"] = rsc

	server.Config.AddDefaultHandler(server.CreateFunctionHandler("/admin/shutdown","ShutDown","GET","RServerCommand"))
	server.Config.AddDefaultHandler(server.CreateFunctionHandler("/admin/loadconfig","LoadConfig","GET","RServerCommand"))
	server.Config.AddDefaultHandler(server.CreateFunctionHandler("/admin/saveconfig","SaveConfig","GET","RServerCommand"))
	server.Config.AddDefaultHandler(server.CreateTemplateHandler("/admin/listcommands","ListCommands","GET","RServerCommand","list","<!DOCTYPE html><html><head><meta charset=\"UTF-8\"><title>Available REST API Calls</title></head><body><h1>Available REST API Calls</h1><h2>Custom</h2>{{range .Handlers}}<div><a href=\"{{ .URL }}\">{{ .URL }} will point to {{ .HTTPMethod }} {{ .FunctionalClass }}.{{ .MethodName }} </a></div>{{else}}<div><strong>None</strong></div>{{end}}<h2>Default</h2>{{range .DefaultHandlers}}<div><a href=\"{{ .URL }}\">{{ .MethodName }}</a></div></div>{{else}}<div><strong>No Handlers</strong></div>{{end}}</body></html>","list.html"))

	for _, handl := range server.Config.GetDefaultHandlers() {
		server.Log.LogDebugf("AddDefaults", "Default Handler: Added %s", handl.MethodName)
	}
}

func SetDefaultFunctionalMap(server *Server.RServer) {
	server.Register("RServerCommand", RServerCommand{Server: server})
}
