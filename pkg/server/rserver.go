package RESTServer

import (
	"context"
	"net/http"
	"reflect"

	Config "github.com/eshu0/RESTServer/pkg/config"
	sl "github.com/eshu0/simplelogger"
	sli "github.com/eshu0/simplelogger/interfaces"
	"github.com/gorilla/mux"
)

// This is the http Server that will host the HTTP requests
var Server *http.Server

type RServer struct {
	Config         Config.IRServerConfig  `json:"-"`
	Log            sli.ISimpleLogger      `json:"-"`
	FunctionalMap  map[string]interface{} `json:"-"`
	ConfigFilePath string                 `json:"-"`
}

func NewRServer(config Config.IRServerConfig) (*RServer) {

	server := RServer{}
	server.Config = config
	server.FunctionalMap = make(map[string]interface{})

	logger := sl.NewApplicationLogger()
	
	// lets open a flie log using the session
	logger.OpenAllChannels()

	server.Log = logger

	return &server
}

func (rs *RServer) Invoke(any interface{}, name string, args ...interface{}) {

	rs.Log.LogDebugf("Invoke", "Method: Looking up %s ", name)

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	val := reflect.ValueOf(any)

	meth := val.MethodByName(name)
	if !meth.IsZero() && !meth.IsNil() {
		meth.Call(inputs)
	} else {
		rs.Log.LogDebugf("Invoke", "Method: %s could not be found ", name)
	}

}

func (rs *RServer) MakeHandler(MethodName string, any interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs.Invoke(any, MethodName, w, r)
	}
}

func (rs *RServer) MapFunctionsToHandlers() *mux.Router {

	r := mux.NewRouter()

	for _, handl := range rs.Config.GetHandlers() {

		funcclass, ok := rs.FunctionalMap[handl.FunctionalClass]

		if ok {
			rs.Log.LogDebugf("MapFunctionsToHandlers", "Handlers: Adding %s", handl.MethodName)
			r.HandleFunc(handl.URL, rs.MakeHandler(handl.MethodName, funcclass)).Methods(handl.HTTPMethod)
		} else {
			rs.Log.LogError("MapFunctionsToHandlers", "Handlers Error FunctionalClass (%s) doesn't have a function mapped", handl.FunctionalClass)
		}
	}

	for _, handl := range rs.Config.GetDefaultHandlers() {

		funcclass, ok := rs.FunctionalMap[handl.FunctionalClass]

		if ok {
			rs.Log.LogDebugf("MapFunctionsToHandlers", "Default Handlers: Adding %s", handl.MethodName)
			r.HandleFunc(handl.URL, rs.MakeHandler(handl.MethodName, funcclass)).Methods(handl.HTTPMethod)
		} else {
			rs.Log.LogError("MapFunctionsToHandlers", "Default Handlers Error FunctionalClass (%s) doesn't have a function mapped", handl.FunctionalClass)
		}

	}

	return r
}

func (rs *RServer) Register(FunctionClass string, data interface{}) {
	rs.FunctionalMap[FunctionClass] = data
}

func (rs *RServer) ShutDown() {
	if Server != nil {
		backg := context.Background()

		if backg != nil {

			if err := Server.Shutdown(backg); err != nil {
				// Error from closing listeners, or context timeout:
				rs.Log.LogDebugf("Shutdown", "HTTP server Shutdown: %v", err)
			}
		} else {
			rs.Log.LogError("Shutdown", "Called but context.Background() was nil")
		}

	} else {
		rs.Log.LogError("Shutdown", "Called but server was nil")
	}

}

func (rs *RServer) ListenAndServe() {
	r := rs.MapFunctionsToHandlers()

	Server = &http.Server{Addr: rs.Config.GetAddress(), Handler: r}

	rs.Log.LogDebugf("ListenAndServe", "Listening on: %s", rs.Config.GetAddress())

	if err := Server.ListenAndServe(); err != http.ErrServerClosed {
		rs.Log.LogErrorf("HTTP server ListenAndServe", "%v", err)
	}
}

func (rs *RServer) SaveConfig() bool {
	ok := rs.Config.Save(rs.ConfigFilePath, rs.Log)
	return ok
}

func (rs *RServer) LoadConfig() bool {
	newconfig, ok := rs.Config.Load(rs.ConfigFilePath, rs.Log)
	if ok {
		rs.Config = newconfig
	}

	return ok
}


