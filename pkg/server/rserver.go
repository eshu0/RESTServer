package RESTServer

import (
	"context"
	"net/http"
	"reflect"
	"html/template"

	Helpers "github.com/eshu0/RESTServer/pkg/helpers"
	Config "github.com/eshu0/RESTServer/pkg/config"

	sl "github.com/eshu0/simplelogger"
	sli "github.com/eshu0/simplelogger/interfaces"
)

// This is the http Server that will host the HTTP requests
var Server *http.Server

type RServer struct {
	Config         		Config.IRServerConfig  `json:"-"`
	Log            		sli.ISimpleLogger      `json:"-"`
	FunctionalMap  		map[string]interface{} `json:"-"`
	ConfigFilePath 		string                 `json:"-"`
	Templates 			*template.Template     `json:"-"`
	RequestHelper 		*Helpers.RequestHelper `json:"-"`
	ResponseHelper 		*Helpers.ResponseHelper `json:"-"`
	NotFoundHandler   	func(w http.ResponseWriter, r *http.Request)
}

func NewRServer(config Config.IRServerConfig) (*RServer) {

	server := RServer{}
	server.Config = config
	server.FunctionalMap = make(map[string]interface{})

	logger := sl.NewApplicationLogger()
	
	// lets open a flie log using the session
	logger.OpenAllChannels()

	server.Log = logger
	server.RequestHelper = Helpers.NewRequestHelper(logger)
	server.ResponseHelper= Helpers.NewResponseHelper(logger)
	return &server
}

func (rs *RServer) Invoke(any interface{}, name string, args ...interface{}) []reflect.Value {

	rs.Log.LogDebugf("Invoke", "Method: Looking up %s ", name)

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		val := reflect.ValueOf(args[i])	
		rs.Log.LogDebugf("Invoke", "ValueOf of arg at [%d] = %v ", i , val)	
		inputs[i] = val
	}
	val := reflect.ValueOf(any)
	rs.Log.LogDebugf("Invoke", "ValueOf %v ", val)
	rs.Log.LogDebugf("Invoke", "Looking up method by %s", name)
	meth := val.MethodByName(name)
	rs.Log.LogDebugf("Invoke", "MethodByName %v ", meth)

	if !meth.IsZero() && !meth.IsNil() {
		return meth.Call(inputs)
	} else {
		rs.Log.LogDebugf("Invoke", "Method: %s could not be found ", name)
	}

	return nil
}

func (rs *RServer) Register(FunctionClass string, data interface{}) {
	rs.FunctionalMap[FunctionClass] = data
}

/// GENERAL OPERATIONS

func (rs *RServer) ShutDown() {
	if Server != nil {
		backg := context.Background()

		if backg != nil {

			if err := Server.Shutdown(backg); err != nil {
				// Error from closing listeners, or context timeout:
				rs.Log.LogErrorf("Shutdown", "HTTP server Shutdown: %v", err)
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

	if rs.NotFoundHandler != nil {
		rs.Log.LogInfo("ListenAndServe", "NotFoundHandler is set")
		r.NotFoundHandler = http.HandlerFunc(rs.NotFoundHandler)
	}

	rs.LoadTemplates()
	Server = &http.Server{Addr: rs.Config.GetAddress(), Handler: r}

	rs.Log.LogInfof("ListenAndServe", "Listening on: %s", rs.Config.GetAddress())

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


