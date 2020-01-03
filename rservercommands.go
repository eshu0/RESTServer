package RESTServer

import (
	"net/http"
)

type RServerCommand struct {
	Server *RServer
}

func (rsc RServerCommand) ShutDown(w http.ResponseWriter, r *http.Request) {
	rsc.Server.ShutDown()
}
