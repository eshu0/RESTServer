package RESTServer

import (
	"net/http"
)

type RServerCommand struct {
	Server *RServer
}

func (rsc RServerCommand) ShutDown(w http.ResponseWriter, r *http.Request) {
	if(rsc.Server != nil){
		rsc.Server.Log.LogDebug("RServerCommand","HTTP server shutdown called")
		rsc.Server.ShutDown()
	}

}
