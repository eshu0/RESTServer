package RESTServer

import (
	"context"
	"html/template"
	"net/http"
	"reflect"

	Config "github.com/eshu0/RESTServer/pkg/config"
	Helpers "github.com/eshu0/RESTServer/pkg/helpers"

	sl "github.com/eshu0/simplelogger/pkg"
	sli "github.com/eshu0/simplelogger/pkg/interfaces"
)

// This is the http Server that will host the HTTP requests
var Server *http.Server

type RServer struct {
	sl.AppLogger
	Config Config.IRServerConfig `json:"-"`
	// This map is designed for the functions were there is no types
	// w http.ResponseWriter, r *http.Request
	RawFunctions map[string]interface{} `json:"-"`
	// This accepts Request.ServerRequest
	TypedMap        map[string]interface{}  `json:"-"`
	ConfigFilePath  string                  `json:"-"`
	Templates       *template.Template      `json:"-"`
	RequestHelper   *Helpers.RequestHelper  `json:"-"`
	ResponseHelper  *Helpers.ResponseHelper `json:"-"`
	NotFoundHandler func(w http.ResponseWriter, r *http.Request)
}

func NewRServer(config Config.IRServerConfig) *RServer {
	return NewRServerCustomLog(config, sl.NewApplicationNowLogger())
}

func NewRServerCustomLog(config Config.IRServerConfig, logger sli.ISimpleLogger) *RServer {
	server := RServer{}
	server.Config = config
	server.Log = logger
	server.RawFunctions = make(map[string]interface{})
	server.TypedMap = make(map[string]interface{})
	server.RequestHelper = Helpers.NewRequestHelper(server.Log)
	server.ResponseHelper = Helpers.NewResponseHelper(server.Log)
	return &server
}

func (rs *RServer) Invoke(any interface{}, name string, args ...interface{}) []reflect.Value {

	rs.LogDebugf("Invoke", "Method: Looking up %s ", name)

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		val := reflect.ValueOf(args[i])
		rs.LogDebugf("Invoke", "ValueOf of arg at [%d] ", i, val)
		rs.LogDebugf("Invoke", "ValueOf of arg at  %v ", val)
		inputs[i] = val
	}

	val := reflect.ValueOf(any)
	rs.LogDebugf("Invoke", "ValueOf %v ", val)
	rs.LogDebugf("Invoke", "Looking up method by %s", name)
	meth := val.MethodByName(name)
	rs.LogDebugf("Invoke", "MethodByName %v", meth)
	/*
		numIn := meth.NumIn() //Count inbound parameters

		for i := 0; i < numIn; i++ {
			inV := meth.In(i)
			in_Kind := inV.Kind() //func
			rs.LogDebugf("Parameter IN: "+strconv.Itoa(i)+"\nKind: %v\nName: %v\n-----------", in_Kind, inV.Name())
		}
	*/
	if meth.IsValid() && !meth.IsNil() {
		return meth.Call(inputs)
	} else {
		rs.LogErrorf("Invoke", "Method: %s could not be found", name)
	}

	return nil
}

//RegisterFunction Registers a func xyz (w http.ResponseWriter, r *http.Request) for processing
func (rs *RServer) RegisterFunction(FunctionClass string, data interface{}) {
	rs.RawFunctions[FunctionClass] = data
}

//Register Registers a func xyz (request Request.ServerRequest) for processing
func (rs *RServer) Register(FunctionClass string, data interface{}) {
	rs.TypedMap[FunctionClass] = data
}

/// GENERAL OPERATIONS

func (rs *RServer) ShutDown() {
	if Server != nil {
		backg := context.Background()

		if backg != nil {

			if err := Server.Shutdown(backg); err != nil {
				// Error from closing listeners, or context timeout:
				rs.LogErrorf("Shutdown", "HTTP server Shutdown: %v", err)
			}
		} else {
			rs.LogError("Shutdown", "Called but context.Background() was nil")
		}

	} else {
		rs.LogError("Shutdown", "Called but server was nil")
	}

}

func (rs *RServer) ListenAndServe() {
	r := rs.MapFunctionsToHandlers()

	if rs.NotFoundHandler != nil {
		rs.LogInfo("ListenAndServe", "NotFoundHandler is set")
		r.NotFoundHandler = http.HandlerFunc(rs.NotFoundHandler)
	}

	rs.LoadTemplates()
	Server = &http.Server{Addr: rs.Config.GetAddress(), Handler: r}

	rs.LogInfof("ListenAndServe", "Listening on: %s", rs.Config.GetAddress())

	if err := Server.ListenAndServe(); err != http.ErrServerClosed {
		rs.LogErrorf("HTTP server ListenAndServe", "%v", err)
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

func DefaultServer(ConfigFilePath *string) (rs *RServer) {

	conf := Config.NewRServerConfig()

	// Create a new REST Server
	server := NewRServer(conf)

	// has a conifg file been provided?
	if ConfigFilePath != nil && len(*ConfigFilePath) > 0 {
		rs.LogDebugf("LoadConfig", "Custom config file path is %s", *ConfigFilePath)

		// load this first
		server.ConfigFilePath = *ConfigFilePath

	} else {
		rs.LogDebugf("LoadConfig", "Default config file path is %s", Config.DefaultFilePath)
		// load this first
		server.ConfigFilePath = Config.DefaultFilePath
	}

	ok := server.LoadConfig()

	// we failed to load the configuration file
	if !ok {
		rs.LogError("LoadConfig ", "failed to load configuration file")
		return
	}

	return server
}
